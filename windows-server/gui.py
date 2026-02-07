import tkinter as tk
from tkinter import ttk, filedialog, messagebox, simpledialog
import socket
from config_manager import ConfigManager
import qrcode
from PIL import Image, ImageTk

class AppGUI:
    def __init__(self, root, config_manager, stop_callback):
        self.root = root
        self.config = config_manager
        self.stop_callback = stop_callback
        
        self.root.title("WebApp Aviator Server")
        self.root.geometry("620x550")
        self.root.minsize(600, 500)
        self.root.protocol("WM_DELETE_WINDOW", self.on_close)

        # Header / Status
        self.status_frame = ttk.Frame(root, padding="10")
        self.status_frame.pack(fill=tk.X)
        
        hostname = socket.gethostname()
        ip = self.get_ip()
        
        ttk.Label(self.status_frame, text="Server Running", font=("Helvetica", 14, "bold")).pack(side=tk.LEFT)
        
        # QR Code Frame
        qr_frame = ttk.Frame(self.status_frame)
        qr_frame.pack(side=tk.RIGHT, padx=10)
        
        self.qr_label = ttk.Label(qr_frame)
        self.qr_label.pack()
        
        # Generate QR
        self.show_qr(f"http://{ip}:8000")

        # URL info
        url_frame = ttk.Frame(self.status_frame)
        url_frame.pack(side=tk.RIGHT, padx=10)
        ttk.Label(url_frame, text=f"Local: http://localhost:8000", foreground="blue").pack(anchor=tk.E)
        ttk.Label(url_frame, text=f"Network: http://{ip}:8000", foreground="green").pack(anchor=tk.E)

        # Buttons (Pack this BEFORE the list frame with side=tk.BOTTOM to ensure visibility)
        self.btn_frame = ttk.Frame(root, padding="10")
        self.btn_frame.pack(side=tk.BOTTOM, fill=tk.X)

        # App List (Pack this with expand=True)
        self.list_frame = ttk.LabelFrame(root, text="Configured Applications", padding="10")
        self.list_frame.pack(fill=tk.BOTH, expand=True, padx=10, pady=5)

        self.tree = ttk.Treeview(self.list_frame, columns=("Name", "Path", "Args"), show="headings", selectmode="browse")
        self.tree.heading("Name", text="Application Name")
        self.tree.heading("Path", text="Executable Path")
        self.tree.heading("Args", text="Arguments")
        self.tree.column("Name", width=150)
        self.tree.column("Path", width=250)
        self.tree.column("Args", width=100)
        self.tree.pack(fill=tk.BOTH, expand=True, side=tk.LEFT)
        
        scrollbar = ttk.Scrollbar(self.list_frame, orient=tk.VERTICAL, command=self.tree.yview)
        self.tree.configure(yscroll=scrollbar.set)
        scrollbar.pack(side=tk.RIGHT, fill=tk.Y)

        # Bind selection event
        self.tree.bind("<<TreeviewSelect>>", self.on_selection_change)

        ttk.Button(self.btn_frame, text="Add Application...", command=self.add_app).pack(side=tk.LEFT, padx=5)
        
        self.btn_edit = ttk.Button(self.btn_frame, text="Edit Args", command=self.edit_args, state="disabled")
        self.btn_edit.pack(side=tk.LEFT, padx=5)
        
        self.btn_remove = ttk.Button(self.btn_frame, text="Remove Selected", command=self.remove_app, state="disabled")
        self.btn_remove.pack(side=tk.LEFT, padx=5)
        
        ttk.Button(self.btn_frame, text="Stop Server", command=self.on_close).pack(side=tk.RIGHT, padx=5)

        self.refresh_list()

    def get_ip(self):
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        try:
            s.connect(('10.255.255.255', 1))
            IP = s.getsockname()[0]
        except Exception:
            IP = '127.0.0.1'
        finally:
            s.close()
        return IP

    def refresh_list(self):
        for item in self.tree.get_children():
            self.tree.delete(item)
        
        for app in self.config.get_apps():
            self.tree.insert("", tk.END, iid=app['id'], values=(app['name'], app['path'], app.get('args', '')))
        
        # Reset buttons on refresh
        self.on_selection_change(None)

    def on_selection_change(self, event):
        selected = self.tree.selection()
        state = "normal" if selected else "disabled"
        self.btn_edit.config(state=state)
        self.btn_remove.config(state=state)

    def add_app(self):
        path = filedialog.askopenfilename(
            title="Select Executable",
            filetypes=[("Executables", "*.exe"), ("All Files", "*.*")]
        )
        if path:
            name = path.split("/")[-1].replace(".exe", "")
            
            # Ask for arguments
            args = simpledialog.askstring("Arguments", f"Optional arguments for {name} (e.g. -bigpicture):")
            if args is None: args = "" # Handle Cancel
            
            self.config.add_app(name, path, args)
            self.refresh_list()

    def edit_args(self):
        selected = self.tree.selection()
        if not selected:
            messagebox.showwarning("Warning", "No application selected")
            return
        
        app_id = selected[0]
        app_item = self.config.get_app_by_id(app_id)
        if not app_item: return

        new_args = simpledialog.askstring("Edit Arguments", 
                                          f"Arguments for {app_item['name']}:", 
                                          initialvalue=app_item.get('args', ''))
        
        if new_args is not None:
            self.config.update_app_args(app_id, new_args)
            self.refresh_list()

    def remove_app(self):
        selected = self.tree.selection()
        if not selected:
            messagebox.showwarning("Warning", "No application selected")
            return
        
        app_id = selected[0]
        self.config.remove_app(app_id)
        self.refresh_list()

    def on_close(self):
        if messagebox.askokcancel("Quit", "Do you want to stop the server?"):
            self.root.destroy()
            self.stop_callback()

    def show_qr(self, data):
        qr = qrcode.QRCode(box_size=4, border=2)
        qr.add_data(data)
        qr.make(fit=True)
        img = qr.make_image(fill="black", back_color="white")
        
        # Convert to Tkinter image
        self.qr_image = ImageTk.PhotoImage(img)
        self.qr_label.config(image=self.qr_image)

def start_gui(config, stop_callback):
    root = tk.Tk()
    app = AppGUI(root, config, stop_callback)
    root.mainloop()
