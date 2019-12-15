import logging
import requests
import time


def main():
    proxies_stack_host = "http://proxies-stack:8080"
    print("Waiting for sleeping scan status...")
    wait_count = 0
    info = requests.get(proxies_stack_host + "/info").json()
    while info["scan_status"] != "sleeping" and not any([info["working"], info["not_working"], info["union"]]):
        time.sleep(2)
        info = requests.get(proxies_stack_host + "/info").json()
        wait_count += 1
        if wait_count > 300:
            raise Exception("Too long waiting!")
    print(info)
    if info["working"] != 2:
        raise Exception("Incorrect amount of the working hosts (%i instead 2)" % info["working"])
    if info["not_working"] != 2:
        raise Exception("Incorrect amount of the not working hosts (%i instead 2)" % info["not_working"])

    print("Checking /working/pop...")
    hosts = []
    resp = requests.get(proxies_stack_host + "/working/pop")
    while resp.status_code == 200:
        hosts.append(resp.json())
        time.sleep(0.5)
        resp = requests.get(proxies_stack_host + "/working/pop")

    info = requests.get(proxies_stack_host + "/info").json()
    if info["not_working"] != 4:
        raise Exception("Incorrect amount of the not working hosts (%i instead 4)" % info["not_working"])

    print("Checking returned hosts")
    must_have_hosts = sorted(["http://foo.com:1080", "http://bar.com:1080"])
    hosts = sorted(hosts)
    if must_have_hosts != hosts:
        raise Exception("Have been received incorrect data (%s instead %s)" % (str(hosts), str(must_have_hosts)))
    

if __name__ == "__main__":
    try:
        main()
    except Exception as ex:
        logging.exception(ex)
        exit(1)
