import socket
from zeroconf import ServiceInfo, Zeroconf

class DiscoveryService:
    def __init__(self, port: int, service_name: str = "AviatorHost"):
        self.port = port
        self.service_name = service_name
        self.zeroconf = Zeroconf()
        self.info = None

    def register(self):
        local_ip = self._get_local_ip()
        
        # Service type: _aviator._tcp.local.
        desc = {'properties': 'v1'}
        
        self.info = ServiceInfo(
            "_aviator._tcp.local.",
            f"{self.service_name}._aviator._tcp.local.",
            addresses=[socket.inet_aton(local_ip)],
            port=self.port,
            properties=desc,
            server=f"{socket.gethostname()}.local."
        )

        print(f"Registering mDNS service: {self.service_name} at {local_ip}:{self.port}")
        self.zeroconf.register_service(self.info)

    def unregister(self):
        if self.info:
            self.zeroconf.unregister_service(self.info)
        self.zeroconf.close()

    def _get_local_ip(self):
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        try:
            # Doesn't reach external, just finds default interface
            s.connect(('10.255.255.255', 1))
            IP = s.getsockname()[0]
        except Exception:
            IP = '127.0.0.1'
        finally:
            s.close()
        return IP
