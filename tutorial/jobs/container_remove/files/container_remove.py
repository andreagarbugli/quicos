import argparse
import sys
import syslog
import docker

import collect_agent


def build_parser():
    parser = argparse.ArgumentParser(description='Docker container removener')
    parser.add_argument('name', type=str, help='Name of docker container to remove')

    return parser


def main(name):
    """Program to stop and remove docker container"""
    try:
        client = docker.from_env()

        container = client.containers.get(name)
        container.stop(timeout=1)
        container.remove()

    except Exception as ex:
        message = 'Error while stopping/removing container {} : {}'.format(name, ex)
        collect_agent.send_log(syslog.LOG_ERR, message)
        sys.exit(message)


if __name__ == '__main__':
    with collect_agent.use_configuration(
            "/opt/openbach/agent/jobs/container_remove/container_remove_rstat_filter.conf"):
        args = build_parser().parse_args()
        main(args.name)
