import time
from selenium import webdriver
from selenium.webdriver.support.wait import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.common import options
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.edge.options import Options

home_page = "http://localhost:5173/"

if __name__ == "__main__":
    edge_options = Options()
    edge_options.add_experimental_option("detach", True)
    driver = webdriver.Edge(options=edge_options)
    driver.get(home_page)

    assert "BringTen" in driver.title

    time.sleep(1)
    wait = WebDriverWait(driver, 2)
    elem = wait.until(EC.element_to_be_clickable((By.ID, "create_game_tab")))
    # elem = driver.find_element(By.ID, "create_game_tab")
    elem.click()

    elem = driver.find_element(By.ID, "create_game_username")
    elem.clear()
    elem.send_keys("Akil")

    elem = driver.find_element(By.ID, "create_game_room_name")
    elem.clear()
    elem.send_keys("Game Grumps")

    elem = driver.find_element(By.ID, "create_game_submit")
    elem.click()

    assert "No results found." not in driver.page_source

    # JOIN GAME PLAYER 2
    driver.switch_to.new_window('tab')
    driver.get(home_page)
    wait = WebDriverWait(driver, 2)
    elem = wait.until(EC.element_to_be_clickable((By.ID, "join_game_tab")))
    elem.click()

    elem = driver.find_element(By.ID, "join_game_username")
    elem.clear()
    elem.send_keys("Des")

    time.sleep(1)
    elem = driver.find_element(By.TAG_NAME, "tbody").find_elements(By.TAG_NAME, "button")[0]
    elem.click()

    assert "No results found." not in driver.page_source

    # JOIN GAME PLAYER 3
    driver.switch_to.new_window('tab')
    driver.get(home_page)
    wait = WebDriverWait(driver, 2)
    elem = wait.until(EC.element_to_be_clickable((By.ID, "join_game_tab")))
    elem.click()

    elem = driver.find_element(By.ID, "join_game_username")
    elem.clear()
    elem.send_keys("Jabari")

    elem = driver.find_element(By.TAG_NAME, "tbody").find_elements(By.TAG_NAME, "button")[0]
    elem.click()

    assert "No results found." not in driver.page_source

    # JOIN GAME PLAYER 4
    driver.switch_to.new_window('tab')
    driver.get(home_page)
    wait = WebDriverWait(driver, 2)
    elem = wait.until(EC.element_to_be_clickable((By.ID, "join_game_tab")))
    elem.click()

    elem = driver.find_element(By.ID, "join_game_username")
    elem.clear()
    elem.send_keys("Momz")

    elem = driver.find_element(By.TAG_NAME, "tbody").find_elements(By.TAG_NAME, "button")[0]
    elem.click()
