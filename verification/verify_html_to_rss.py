import re
from playwright.sync_api import sync_playwright, Page, expect

def run(playwright):
    browser = playwright.chromium.launch(headless=True)
    context = browser.new_context()
    page = context.new_page()

    # Mock API responses
    page.route("**/api/admin/user/info", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"code":0,"msg":"success","data":{"id":1,"username":"admin","permissions":["*"],"avatar":"https://avatars.githubusercontent.com/u/1?v=4"}}'
    ))

    page.route("**/api/admin/menu/list", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"code":0,"msg":"success","data":[]}'
    ))

    page.route("**/api/admin/tools/fetch", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"code":0,"msg":"success","data":"<html><body><div class=\\"article-card\\"><h1>Title</h1></div></body></html>"}'
    ))

    page.add_init_script("""
        localStorage.setItem('token', 'mock-token');
        localStorage.setItem('arco-locale', 'en-US');
    """)

    page.goto("http://localhost:5173/tools/html-to-rss")

    # Step 1
    page.get_by_placeholder("https://example.com/blog").fill("https://example.com")
    page.get_by_role("button", name="Fetch & Next").click()

    # Wait for Step 2
    expect(page.get_by_text("Page Preview")).to_be_visible()

    # Target the Item Selector Pick button (first suffix button)
    # We wait for it to be visible first.
    pick_btn = page.locator(".arco-input-suffix button").first
    expect(pick_btn).to_be_visible()

    # Verification 1: Auto-advance should activate Picking
    print("Verifying initial active state...")
    expect(pick_btn).to_have_text("Picking...")

    # Take screenshot of Active State
    pick_btn.screenshot(path="verification/btn_active.png")
    page.screenshot(path="verification/state_active.png")

    # Verification 2: Click to toggle OFF
    print("Verifying toggle off...")
    pick_btn.click()
    expect(pick_btn).to_have_text("Pick")

    # Verification 3: Click to toggle ON
    print("Verifying toggle on...")
    pick_btn.click()
    expect(pick_btn).to_have_text("Picking...")

    print("Verification Successful!")
    browser.close()

if __name__ == "__main__":
    with sync_playwright() as playwright:
        run(playwright)
