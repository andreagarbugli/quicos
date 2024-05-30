import argparse
import sys
import syslog

import socket
import json

import docker

import collect_agent

def build_parser():
    parser = argparse.ArgumentParser(description='Docker container runner')
    parser.add_argument('docker_image', type=str, help='Name of docker image to run')
    parser.add_argument('-c', '--command', type=str, help='Command to run on the container', default="")
    parser.add_argument('-n', '--name', type=str, help='Container\'s name', default="")
    parser.add_argument('-p', '--ports', type=str, help='List of comma separated ports to bind', default="")
    parser.add_argument('-m', '--metrics', type=str, help='Container port where to listen to metrics', default="")

    return parser


def main(image, command, name, ports, metrics):
    """Program to run docker container"""
    try:
        client = docker.from_env()
        list_port = {}

        if not command:
            command = None

        if ports:
            ports = ports.split(',')

            for port in ports:
                parts = port.split(':')

                if len(parts) == 2:
                    key = f'{parts[1]}/tcp'
                    value = int(parts[0])
                    list_port[key] = value

        if not name:
            client.containers.run(
                image=image,
                command=command,
                detach=True,
                ports=list_port
            )
        else:
            client.containers.run(
                image=image,
                command=command,
                detach=True,
                name=name,
                ports=list_port
            )

        # Create a socket object
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

        try:
            # Connect to the server
            sock.connect(('localhost', int(metrics)))

            while True:
                # Receive data from the socket
                data = sock.recv(1024)

                # Decode the received data
                received_data = data.decode()

                try:
                    # Unmarshal the JSON data
                    json_data = json.loads(received_data)
                    timestamp = int(json_data['timestamp'])
                    name = json_data['name']
                    value = float(json_data['value'])

                    print(timestamp, name, value)

                    if name == 'ping':
                        collect_agent.send_stat(timestamp, ping=value)
                    else:
                        collect_agent.send_stat(timestamp, stat1=value)

                except json.JSONDecodeError as e:
                    message = 'Error while decoding json metric {}'.format(e)
                    collect_agent.send_log(syslog.LOG_ERR, message)
                    sys.exit(message)

        except ConnectionRefusedError:
            message = 'Connection to metrics port refused.'
            collect_agent.send_log(syslog.LOG_ERR, message)
            sys.exit(message)
        except Exception as ex:
            message = 'Error with socket: {}'.format(ex)
            collect_agent.send_log(syslog.LOG_ERR, message)
            sys.exit(message)

        finally:
            # Close the socket
            sock.close()

    except Exception as ex:
        message = 'Error while executing container {} : {}'.format(image, ex)
        collect_agent.send_log(syslog.LOG_ERR, message)
        sys.exit(message)


if __name__ == '__main__':
    with collect_agent.use_configuration("/opt/openbach/agent/jobs/container_run_server"
                                         "/container_run_server_rstat_filter.conf"):
        args = build_parser().parse_args()
        main(args.docker_image, args.command, args.name, args.ports, args.metrics)
