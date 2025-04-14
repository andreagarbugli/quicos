from auditorium_scripts.install_jobs import InstallJobs
from auditorium_scripts.uninstall_jobs import UninstallJobs

def main():
    install_jobs = InstallJobs()
    install_jobs.parse('-j iperf -a 192.168.50.11 -j iperf -a 192.168.50.10'.split())
    install_jobs.execute()

    # uninstall_jobs = UninstallJobs()
    # uninstall_jobs.parse('-j iperf -a 192.168.50.11 -j iperf -a 192.168.50.10'.split())
    # uninstall_jobs.execute()


if __name__ == "__main__":
    main()