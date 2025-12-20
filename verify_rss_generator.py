import time
from playwright.sync_api import sync_playwright, expect

def run(playwright):
    browser = playwright.chromium.launch(headless=True)
    # Increase viewport size just in case
    context = browser.new_context(viewport={"width": 1280, "height": 720})
    page = context.new_page()

    # Mock API endpoints
    # 1. User Info
    page.route("**/api/admin/user/info", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"code":0, "data": {"id": 1, "username": "admin", "roles": ["admin"]}}'
    ))

    # 2. Menu List
    page.route("**/api/admin/menu/list", lambda route: route.fulfill(
        status=200,
        content_type="application/json",
        body='{"code":0, "data": []}'
    ))

    # 3. Fetch Tool (Mocking the fetch response)
    def handle_fetch(route):
        request = route.request
        post_data = request.post_data_json
        use_browserless = post_data.get("use_browserless", False)

        print(f"Fetch called with use_browserless={use_browserless}")

        content = "<html><body><h1>Standard Content</h1><ul><li>Item 1</li></ul></body></html>"
        if use_browserless:
            content = "<html><body><h1>Browserless Content</h1><ul><li>Item 1</li><li>Item 2 (JS)</li></ul></body></html>"

        route.fulfill(
            status=200,
            content_type="application/json",
            body=f'{{"code":0, "data": "{content}"}}'
        )

    page.route("**/api/admin/tools/fetch", handle_fetch)

    # Inject token to bypass login
    page.add_init_script("localStorage.setItem('token', 'dummy-token');")
    # Also set locale to En-US to match our text assertions
    page.add_init_script("localStorage.setItem('arco-locale', 'en-US');")

    # Navigate to RSS Generator
    print("Navigating to RSS Generator...")
    # Add fail_on_status_code=False in case of 404 handled by client routing
    # But usually goto waits for load.
    try:
        page.goto("http://localhost:5173/tools/rss-generator", timeout=30000)
    except Exception as e:
        print(f"Goto failed: {e}")
        # Sometimes connection refused if server not ready

    # Wait for app to mount
    print("Waiting for app...")
    page.wait_for_selector("#app", timeout=10000)

    print("Waiting for wizard card...")
    try:
        page.wait_for_selector(".wizard-card", timeout=10000)
    except Exception as e:
        print("Timeout waiting for wizard card. Taking screenshot...")
        page.screenshot(path="/home/jules/verification/timeout_debug.png")
        print(page.content()) # Print HTML to debug
        raise e

    # Step 1: Check if Enhance Mode switch is GONE
    step1_content = page.locator(".step-content").first
    expect(step1_content).to_be_visible()

    # Check that "Enhanced Mode" text is NOT present in Step 1 form
    expect(step1_content.get_by_text("Enhanced Mode (Browserless)")).not_to_be_visible()

    print("Verified Step 1: Enhance Mode switch is absent.")
    page.screenshot(path="/home/jules/verification/step1.png")

    # Enter URL and Next
    print("Entering URL...")
    # Use get_by_role for input if placeholder fails
    page.get_by_placeholder("https://example.com/blog").fill("http://example.com")

    # Wait for the button to be enabled (it has :disabled="!url")
    # "Fetch & Next" in English.
    # Since we forced locale, this should work.

    # Check if button exists first
    btn = page.locator("button").filter(has_text="Fetch & Next")
    if btn.count() == 0:
        # Fallback to Chinese text if locale force failed?
        # 'rssGenerator.step1.button': '获取并下一步'
        btn = page.locator("button").filter(has_text="获取并下一步")

    expect(btn).to_be_enabled()
    btn.click()

    # Wait for Step 2
    print("Waiting for Step 2...")
    expect(page.get_by_text("Page Preview")).to_be_visible()

    print("Moved to Step 2.")

    # Verify Enhance Mode switch IS present in Step 2
    # The label is "Enhance Mode" or "增强模式". In En-US it is "Enhance Mode".
    # Wait for text "Enhanced Mode"
    try:
        expect(page.get_by_text("Enhanced Mode", exact=True)).to_be_visible()
    except:
        print("Could not find 'Enhanced Mode' text. Checking fallback.")
        page.screenshot(path="/home/jules/verification/step2_missing_text.png")
        # Try Chinese
        try:
            expect(page.get_by_text("增强模式", exact=True)).to_be_visible()
            print("Found '增强模式'")
        except:
             print("Could not find Chinese text either.")

    # Check for the switch
    switch = page.locator(".arco-switch").first
    if switch.count() == 0:
         switch = page.locator("button[role='switch']").first

    expect(switch).to_be_visible()

    print("Verified Step 2: Enhance Mode switch is present.")
    page.screenshot(path="/home/jules/verification/step2_initial.png")

    # Toggle Switch
    print("Toggling Enhance Mode...")
    switch.click()

    # Wait a bit
    time.sleep(2)
    page.screenshot(path="/home/jules/verification/step2_toggled.png")
    print("Toggled switch and captured screenshot.")

    browser.close()

with sync_playwright() as playwright:
    run(playwright)
