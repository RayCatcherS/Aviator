import subprocess
import os
import threading

def run_executable(path: str, args: str = ""):
    """
    Runs the executable at the given path with optional arguments.
    """
    if not os.path.exists(path):
        raise FileNotFoundError(f"Executable not found at: {path}")

    # Prepare command
    cmd = [path]
    if args:
        # Split args string into list roughly (simple split)
        # For complex args with quotes, shlex might be better, but keep it simple for now
        import shlex
        try:
            cmd.extend(shlex.split(args))
        except:
             cmd.extend(args.split())

    # Using subprocess.Popen is generally more robust for args than startfile
    try:
        # cwd=os.path.dirname(path) is often important for games/apps to find their assets
        subprocess.Popen(cmd, cwd=os.path.dirname(path))
    except Exception as e:
        print(f"Error launching: {e}")
        # Fallback if Popen fails? 
        # os.startfile generally doesn't take args easily
        pass
