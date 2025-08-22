from selenium import webdriver
from selenium.webdriver.chrome.options import Options
import os, shutil, tempfile, atexit
import time

SLEEP = 5
COUNT = 10
URL = "https://kevin-fares-prospective-fame.trycloudflare.com/"

SCRIPT = """
var buttons = document.getElementsByTagName('button');
buttons[0].click()
"""

def access():
    options = Options()
    options.add_experimental_option("excludeSwitches", ["enable-logging"])

    options.add_experimental_option(
         "prefs", {"profile.default_content_setting_values.notifications": 1}
    )

    with webdriver.Chrome(options=options) as driver:
        driver.get(URL)
        time.sleep(SLEEP)
        driver.execute_script(SCRIPT)
        time.sleep(SLEEP)

def main():
    for i in range(COUNT):
        print(i)
        access()


if __name__ == "__main__":
    main()