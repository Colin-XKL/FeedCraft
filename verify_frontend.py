from playwright.sync_api import sync_playwright, expect
import time

def test_custom_recipe_multi_select(page):
    # Mock APIs
    page.route("**/api/admin/user/info", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"code":0,"msg":"ok","data":{"name":"admin","role":"admin"}}'
    ))
    page.route("**/api/admin/menu/list", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"code":0,"msg":"ok","data":[]}'
    ))
    page.route("**/api/admin/recipes", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"code":0,"msg":"ok","data":[]}'
    ))

    # Mock Craft APIs
    page.route("**/api/admin/sys-craft-atoms", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"code":0,"msg":"ok","data":[{"name":"fulltext","description":"Fulltext extraction"},{"name":"proxy","description":"Proxy feed"}]}'
    ))
    page.route("**/api/admin/craft-atoms", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"code":0,"msg":"ok","data":[{"name":"my-atom","description":"User atom"}]}'
    ))
    page.route("**/api/admin/craft-flows", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"code":0,"msg":"ok","data":[{"name":"my-flow","description":"User flow"}]}'
    ))
    page.route("**/api/admin/craft-templates", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"code":0,"msg":"ok","data":[{"name":"fulltext","description":"Fulltext extraction"},{"name":"proxy","description":"Proxy feed"}]}'
    ))

    # Set local storage for auth and locale
    page.add_init_script("""
        localStorage.setItem('token', 'mock-token');
        localStorage.setItem('arco-locale', 'en-US');
    """)

    # Go to Custom Recipe page
    page.goto("http://localhost:5173/dashboard/custom_recipe")

    # Wait for page load
    expect(page.get_by_role("button", name="Create", exact=True)).to_be_visible()

    # Click Create
    page.get_by_role("button", name="Create", exact=True).click()

    # Expect Modal
    # We look for the modal title which changes depending on edit/create
    expect(page.get_by_text("Create Recipe")).to_be_visible()

    # Find Craft Selector trigger
    page.locator(".trigger-input").click()

    # Expect Selection Modal
    # The selector modal title matches placeholder.
    # In my code: :title="placeholder || t('feedCompare.selectCraftFlow.placeholder')"
    # In custom_recipe: :placeholder="t('customRecipe.form.placeholder.craft')"
    # I assume "Select Craft" or something.
    # I'll just wait for the dialog to be visible.
    time.sleep(1) # Wait for animation

    # Select "fulltext" (System)
    page.get_by_text("fulltext", exact=True).first.click()

    # Select "proxy" (System)
    page.get_by_text("proxy", exact=True).first.click()

    # Click Confirm (OK) in the selection modal
    # Arco Design Modal OK button
    page.get_by_role("button", name="OK").click()

    # Expect tags in the main modal
    expect(page.locator(".trigger-input")).to_contain_text("fulltext")
    expect(page.locator(".trigger-input")).to_contain_text("proxy")

    # Screenshot
    page.screenshot(path="/home/jules/verification/custom_recipe_multi_select.png")

if __name__ == "__main__":
    with sync_playwright() as p:
        browser = p.chromium.launch(headless=True)
        page = browser.new_page()
        try:
            test_custom_recipe_multi_select(page)
        except Exception as e:
            print(f"Error: {e}")
            page.screenshot(path="/home/jules/verification/error.png")
            raise e
        finally:
            browser.close()
