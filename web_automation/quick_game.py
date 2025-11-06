import time
from selenium import webdriver
from selenium.webdriver.support.wait import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common import options, service
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.edge.options import Options

home_page = "http://localhost:5173/"

def create_game(driver):
    driver.get(home_page)
    assert "BringTen" in driver.title

    time.sleep(0.5)
    elem = driver.find_element(By.ID, "create_game_tab")
    elem.click()

    driver.find_element(By.ID, "create_game_username").send_keys("Akil")
    driver.find_element(By.ID, "create_game_room_name").send_keys("Game Grumps")
    driver.find_element(By.ID, "create_game_submit").click()

    assert "No results found." not in driver.page_source


def join_game(driver, username):
    driver.get(home_page)
    wait = WebDriverWait(driver, 5)

    time.sleep(0.5)
    driver.find_element(By.ID, "join_game_tab").click()
    driver.find_element(By.ID, "join_game_username").send_keys(username)

    # Wait for the list of games to appear and click the first "Join" button
    time.sleep(1)  # Let game list load
    buttons = driver.find_elements(By.CLASS_NAME, "joinGameBtn")
    if buttons:
        buttons[0].click()
    else:
        raise Exception("No join buttons found!")

    assert "No results found." not in driver.page_source


if __name__ == "__main__":
    edge_options = Options()
    edge_options.add_experimental_option("detach", True)  # Set to False so browser closes on failure

    service = webdriver.EdgeService(executable_path="./edgedriver_mac64_m1/msedgedriver")

    driver = webdriver.Edge(service=service, options=edge_options)

    # CREATE GAME
    create_game(driver)

    # JOIN GAME IN NEW TABS
    usernames = ["Des", "Jabari", "Momz"]
    for name in usernames:
        driver.switch_to.new_window('tab')
        join_game(driver, name)

    print("✅ All players successfully joined the game.")

