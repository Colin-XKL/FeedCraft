from playwright.sync_api import sync_playwright, expect

def run(playwright):
    browser = playwright.chromium.launch(headless=True)
    context = browser.new_context()
    page = context.new_page()

    page.on("console", lambda msg: print(f"Console: {msg.text}"))

    # Mock API
    def handle_fetch(route):
        route.fulfill(
            status=200,
            content_type="application/json",
            body='{"code": 0, "msg": "ok", "data": "<html><head><title>Test Feed Title</title></head><body><div class=\\"item\\"><h1>Item 1</h1></div></body></html>"}'
        )
    page.route("**/api/admin/tools/fetch", handle_fetch)

    def handle_parse(route):
        route.fulfill(
            status=200,
            content_type="application/json",
            body='{"code": 0, "msg": "ok", "data": [{"title": "Item 1", "link": "http://example.com/1", "date": "2023-01-01", "description": "Desc 1"}]}'
        )
    page.route("**/api/admin/tools/parse", handle_parse)

    print("Navigating...")
    try:
        page.goto("http://localhost:5173/worktable/html-to-rss")
        page.wait_for_timeout(3000)

        print("Step 1...")
        page.get_by_placeholder("https://example.com/blog").fill("http://example.com")
        page.locator(".step-content button.arco-btn-primary").first.click()

        page.wait_for_selector("text=页面预览", timeout=10000)

        print("Step 2...")
        # Step 2 is the second .step-content
        step2 = page.locator(".step-content").nth(1)
        step2.get_by_placeholder(".article-card").fill(".item")

        step2.locator("button:has-text('运行预览')").click()
        page.wait_for_selector("text=Item 1", timeout=5000)

        step2.locator("button:has-text('下一步')").click()

        print("Step 3...")
        # Step 3 is the third .step-content
        step3 = page.locator(".step-content").nth(2)
        expect(step3.get_by_placeholder("例如：我的技术博客 RSS")).to_have_value("Test Feed Title")
        step3.locator("button:has-text('下一步')").click()

        print("Step 4...")
        # Step 4 is the fourth .step-content
        step4 = page.locator(".step-content").nth(3)
        recipe_id = step4.get_by_placeholder("my-recipe-id")
        expect(recipe_id).to_have_value("test-feed-title")

        print("Testing Regeneration...")
        step4.locator("button:has-text('返回')").click()

        # Back to Step 3
        step3.get_by_placeholder("例如：我的技术博客 RSS").fill("Updated Title 123")
        step3.locator("button:has-text('下一步')").click()

        expect(recipe_id).to_have_value("test-feed-title")

        step4.locator(".arco-input-append button").click()

        expect(recipe_id).to_have_value("updated-title-123")

        print("Taking success screenshot...")
        page.screenshot(path="verification/verification.png")

    except Exception as e:
        print(f"Error: {e}")
        page.screenshot(path="verification/error.png")

    finally:
        context.close()
        browser.close()

with sync_playwright() as p:
    run(p)
