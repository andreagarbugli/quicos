import argparse
import sys
import syslog
import docker

import collect_agent


def build_parser():
    parser = argparse.ArgumentParser(description="Docker container installer")
    parser.add_argument("docker_repo", type=str, help="Name of docker image")
    parser.add_argument(
        "-t", "--tag", type=str, help="Name of docker image tag", default="latest"
    )
    parser.add_argument(
        "-r",
        "--registry",
        type=str,
        help="Name of the docker private registry",
        default="",
    )
    parser.add_argument(
        "-u",
        "--username",
        type=str,
        help="Username to access private registry",
        default="",
    )
    parser.add_argument(
        "-p",
        "--password",
        type=str,
        help="Password to access private registry",
        default="",
    )

    return parser


def main(repository, tag, registry, username, password):
    """Program to install docker image"""
    try:
        client = docker.from_env()

        if registry and username and password:
            auth_config = {"username": username, "password": password}

            client.images.pull(
                repository=registry + repository, tag=tag, auth_config=auth_config
            )

        else:
            client.images.pull(repository=repository, tag=tag)
    except Exception as ex:
        message = "Error while pulling image {}{} : {}".format(repository, tag, ex)
        collect_agent.send_log(syslog.LOG_ERR, message)
        sys.exit(message)


if __name__ == "__main__":
    with collect_agent.use_configuration(
        "/opt/openbach/agent/jobs/container_installer/container_installer_rstat_filter.conf"
    ):
        args = build_parser().parse_args()
        main(args.docker_repo, args.tag, args.registry, args.username, args.password)
