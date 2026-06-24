package email

import "fmt"

const fontStack = "'Plus Jakarta Sans', -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif"

// emailWrapper wraps content in the WashPoint branded email shell.
func emailWrapper(title, statusBg, statusFg, statusIcon, statusLabel, body string) string {
	iconHTML := ""
	if statusIcon != "" {
		iconHTML = `<span style="font-size:18px;margin-right:8px;">` + statusIcon + `</span>`
	}

	return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8"/>
<meta name="viewport" content="width=device-width,initial-scale=1"/>
<title>` + title + `</title>
</head>
<body style="margin:0;padding:0;background:#eceae7;font-family:` + fontStack + `;">
  <div style="padding:32px 16px 48px;">

    <!-- Email card -->
    <div style="max-width:640px;margin:0 auto;background:#ffffff;border-radius:6px;overflow:hidden;box-shadow:0 4px 32px rgba(0,0,0,0.10);">

      <!-- Header -->
      <table width="100%" cellpadding="0" cellspacing="0" style="background:linear-gradient(135deg,#F26F21 0%,#D85D14 100%);">
        <tr>
          <td align="center" style="padding:38px 40px 34px;">
            <table cellpadding="0" cellspacing="0">
              <tr>
                <td align="center" style="padding-bottom:6px;">
                  <span style="font-size:26px;font-weight:800;color:#ffffff;letter-spacing:-0.5px;font-family:` + fontStack + `;">💧 WashPoint</span>
                </td>
              </tr>
              <tr>
                <td align="center">
                  <span style="font-size:11px;font-weight:700;letter-spacing:2px;color:rgba(255,255,255,0.75);font-family:` + fontStack + `;">PROFESSIONAL LAUNDRY SERVICES</span>
                </td>
              </tr>
            </table>
          </td>
        </tr>
      </table>

      <!-- Status bar -->
      <table width="100%" cellpadding="0" cellspacing="0" style="background:` + statusBg + `;">
        <tr>
          <td align="center" style="padding:13px 40px;">
            ` + iconHTML + `<span style="font-size:13px;font-weight:700;color:` + statusFg + `;font-family:` + fontStack + `;">` + statusLabel + `</span>
          </td>
        </tr>
      </table>

      <!-- Body -->
      <div style="padding:36px 40px 28px;">
        ` + body + `
      </div>

      <!-- Divider -->
      <div style="height:1px;background:#f1eeea;margin:0 40px;"></div>

      <!-- Footer -->
      <div style="padding:24px 40px;text-align:center;background:#faf8f6;">
        <div style="font-size:13.5px;font-weight:700;color:#1a1f2e;margin-bottom:6px;">WashPoint Laundry Services</div>
        <div style="font-size:12px;color:#aab0ba;line-height:1.7;">
          Mon–Fri: 07:00–19:00 &nbsp;|&nbsp; Sat: 08:00–17:00 &nbsp;|&nbsp; Sun: Closed
        </div>
        <div style="margin-top:10px;font-size:12px;color:#aab0ba;">
          <span>📞 +260 97X XXX XXX</span>&nbsp;&nbsp;·&nbsp;&nbsp;<span>📍 Lusaka, Zambia</span>
        </div>
      </div>

    </div>

    <div style="max-width:640px;margin:18px auto 0;text-align:center;font-size:11.5px;color:#b0b7c0;font-weight:500;">
      This email was sent by WashPoint Laundry Services · Lusaka, Zambia
    </div>
  </div>
</body>
</html>`
}

// orderPill renders the order reference pill (order number / service / date label).
func orderPill(orderNumber, service, dateLabel, dateValue string) string {
	return fmt.Sprintf(`
      <div style="display:flex;align-items:center;gap:16px;background:#faf8f6;border:1px solid #ece8e3;border-radius:12px;padding:14px 18px;margin-bottom:28px;">
        <div style="flex:1;">
          <div style="font-size:11px;font-weight:700;letter-spacing:0.8px;color:#aab0ba;">ORDER NUMBER</div>
          <div style="font-size:19px;font-weight:800;color:#1a1f2e;letter-spacing:-0.3px;">%s</div>
        </div>
        <div style="flex:1;border-left:1px solid #ece8e3;padding-left:16px;">
          <div style="font-size:11px;font-weight:700;letter-spacing:0.8px;color:#aab0ba;">SERVICE</div>
          <div style="font-size:14px;font-weight:700;color:#1a1f2e;">%s</div>
        </div>
        <div style="flex:1;border-left:1px solid #ece8e3;padding-left:16px;">
          <div style="font-size:11px;font-weight:700;letter-spacing:0.8px;color:#aab0ba;">%s</div>
          <div style="font-size:14px;font-weight:700;color:#1a1f2e;">%s</div>
        </div>
      </div>`, orderNumber, service, dateLabel, dateValue)
}

// itemsTable renders the items checklist table.
func itemsTable(items []struct {
	Name string
	Qty  int
}) string {
	rows := ""
	total := 0
	for _, it := range items {
		letter := "•"
		if len(it.Name) > 0 {
			letter = string([]rune(it.Name)[0])
		}
		rows += fmt.Sprintf(`
        <div style="display:flex;align-items:center;justify-content:space-between;padding:14px 18px;border-top:1px solid #f1eeea;">
          <div style="display:flex;align-items:center;gap:11px;">
            <div style="width:32px;height:32px;border-radius:9px;background:#FEF1E9;color:#D85D14;display:inline-flex;align-items:center;justify-content:center;font-weight:800;font-size:13px;flex:none;line-height:1;text-align:center;">%s</div>
            <span style="font-size:14px;font-weight:600;color:#222831;">%s</span>
          </div>
          <span style="font-size:14px;font-weight:700;color:#5b6472;">× %d</span>
        </div>`, letter, it.Name, it.Qty)
		total += it.Qty
	}

	return fmt.Sprintf(`
      <div style="border:1px solid #ece8e3;border-radius:12px;overflow:hidden;margin-bottom:24px;">
        <div style="display:flex;align-items:center;justify-content:space-between;padding:12px 18px;background:#f6f4f1;">
          <span style="font-size:11px;font-weight:800;letter-spacing:1px;color:#aab0ba;">ITEM</span>
          <span style="font-size:11px;font-weight:800;letter-spacing:1px;color:#aab0ba;">QTY</span>
        </div>
        %s
        <div style="display:flex;align-items:center;justify-content:space-between;padding:13px 18px;background:#f6f4f1;border-top:1px solid #ece8e3;">
          <span style="font-size:13px;font-weight:700;color:#5b6472;">Total items</span>
          <span style="font-size:14px;font-weight:800;color:#F26F21;">%d&nbsp;pcs</span>
        </div>
      </div>`, rows, total)
}

// collectionBox renders the collection details block.
func collectionBox() string {
	return `
      <div style="border:1px solid #ece8e3;border-radius:12px;padding:18px 20px;margin-bottom:28px;">
        <div style="margin-bottom:14px;">
          <span style="font-size:13.5px;font-weight:800;color:#1a1f2e;">📍 Collection Details</span>
        </div>
        <div style="font-size:13.5px;color:#5b6472;line-height:1.5;">
          <div style="margin-bottom:9px;"><span style="font-weight:700;color:#1a1f2e;">Location:</span> WashPoint, Lusaka, Zambia</div>
          <div style="margin-bottom:9px;"><span style="font-weight:700;color:#1a1f2e;">Hours:</span> Mon–Fri 07:00–19:00 | Sat 08:00–17:00</div>
          <div><span style="font-weight:700;color:#1a1f2e;">Phone:</span> +260 97X XXX XXX</div>
        </div>
      </div>`
}

// CustomerWelcomeTemplate returns the branded welcome email for a newly registered customer.
func CustomerWelcomeTemplate(customerName string) string {
	body := fmt.Sprintf(`
      <h1 style="font-size:22px;font-weight:800;color:#1a1f2e;line-height:1.3;margin:0 0 12px;letter-spacing:-0.3px;">Welcome to WashPoint, %s!</h1>
      <p style="font-size:14.5px;color:#5b6472;line-height:1.7;margin:0 0 28px;">Your account has been created and you're all set to experience professional laundry care that saves you time and keeps your clothes looking their best.</p>

      <!-- Services -->
      <div style="border:1px solid #ece8e3;border-radius:12px;overflow:hidden;margin-bottom:24px;">
        <div style="padding:12px 18px;background:#f6f4f1;">
          <span style="font-size:11px;font-weight:800;letter-spacing:1px;color:#aab0ba;">OUR SERVICES</span>
        </div>
        <div style="font-size:14px;font-weight:600;color:#222831;">
          <div style="padding:11px 18px;border-top:1px solid #f1eeea;">👔 &nbsp;Wash &amp; Fold</div>
          <div style="padding:11px 18px;border-top:1px solid #f1eeea;">🧥 &nbsp;Dry Cleaning</div>
          <div style="padding:11px 18px;border-top:1px solid #f1eeea;">👗 &nbsp;Ironing &amp; Pressing</div>
          <div style="padding:11px 18px;border-top:1px solid #f1eeea;">🛏️ &nbsp;Bedding &amp; Linen</div>
          <div style="padding:11px 18px;border-top:1px solid #f1eeea;">🧣 &nbsp;Delicate Care</div>
          <div style="padding:11px 18px;border-top:1px solid #f1eeea;">👔 &nbsp;Suit &amp; Formal Wear</div>
        </div>
      </div>

      <!-- How it works -->
      <div style="border:1px solid #ece8e3;border-radius:12px;padding:18px 20px;margin-bottom:28px;">
        <div style="font-size:11px;font-weight:800;letter-spacing:1px;color:#aab0ba;margin-bottom:14px;">HOW IT WORKS</div>
        <div style="font-size:13.5px;color:#5b6472;line-height:1.5;">
          <div style="margin-bottom:10px;"><span style="font-weight:700;color:#1a1f2e;">1.</span> &nbsp;Drop off your laundry at our shop</div>
          <div style="margin-bottom:10px;"><span style="font-weight:700;color:#1a1f2e;">2.</span> &nbsp;We sort, clean, and care for each item</div>
          <div style="margin-bottom:10px;"><span style="font-weight:700;color:#1a1f2e;">3.</span> &nbsp;Get notified when your order is ready</div>
          <div><span style="font-weight:700;color:#1a1f2e;">4.</span> &nbsp;Pick up fresh, clean clothes — simple!</div>
        </div>
      </div>

      `+collectionBox()+`

      <p style="font-size:13.5px;color:#8a8f98;line-height:1.7;margin:0 0 20px;">If you have any questions or special requests, don't hesitate to reach out to us at the shop. We're always happy to help!</p>
      <p style="font-size:14px;color:#5b6472;margin:0;">See you soon,<br><strong style="color:#F26F21;">The WashPoint Team</strong></p>`,
		customerName,
	)

	return emailWrapper(
		"Welcome to WashPoint!",
		"#eef2ff", "#4154d6", "👋",
		"Welcome — your account is ready!",
		body,
	)
}

// OrderReadyTemplate returns the branded pickup-ready email.
func OrderReadyTemplate(customerName, orderNumber, serviceType string, items []struct {
	Name string
	Qty  int
}) string {
	body := fmt.Sprintf(`
      <h1 style="font-size:22px;font-weight:800;color:#1a1f2e;line-height:1.3;margin:0 0 12px;letter-spacing:-0.3px;">Your laundry is ready, %s!</h1>
      <p style="font-size:14.5px;color:#5b6472;line-height:1.7;margin:0 0 28px;">Great news — your %s order <strong style="color:#1a1f2e;">%s</strong> has been processed and is now ready for collection. Please come pick it up at your earliest convenience.</p>

      %s

      %s

      <!-- Amount due -->
      <div style="display:flex;align-items:center;justify-content:space-between;background:#FEF1E9;border:1px solid #FBD9C2;border-radius:12px;padding:16px 20px;margin-bottom:24px;">
        <div>
          <div style="font-size:11px;font-weight:700;letter-spacing:0.8px;color:#D85D14;">PAYMENT</div>
          <div style="font-size:14px;font-weight:700;color:#D85D14;">Balance due on pickup</div>
        </div>
      </div>

      %s

      <p style="font-size:13.5px;color:#8a8f98;line-height:1.7;margin:0 0 20px;">Please bring this email or your order number <strong style="color:#1a1f2e;">%s</strong> when you come to collect. If you have any questions, don't hesitate to call us.</p>
      <p style="font-size:14px;color:#5b6472;margin:0;">See you soon,<br><strong style="color:#F26F21;">The WashPoint Team</strong></p>`,
		customerName,
		serviceType, orderNumber,
		orderPill(orderNumber, serviceType, "READY SINCE", "Today"),
		itemsTable(items),
		collectionBox(),
		orderNumber,
	)

	return emailWrapper(
		fmt.Sprintf("Your laundry is ready for pickup — %s", orderNumber),
		"#e8f7ee", "#1f9d57", "",
		"Your laundry is ready for pickup!",
		body,
	)
}
