Hitbox
======

#### A lightweight server application to alert you to your site's visitors via Slack.

Does your website receive so few visitors that you'd like to be alerted to each and every one? Perfect! And sorry.

Hitbox listens for GET requests, generated by a snippet of client-side Javascript, that supply basic page and referrer information. Hitbox then uses the client's IP address to query a GeoIP service (by default, [freegeoip.net](http://freegeoip.net)) to provide a best-guess of the user's city, region, and country.

## Installation
Building Hitbox is as easy as building any other Golang application.  If you're planning on deploying Hitbox on server, ensure that you set your `GOOS` and `GOARCH` environment variables to match the operating system and architecture of the target platform.  Hitbox has no dependencies.

```
export GOOS=<target os>
export GOARCH=<target architecture>

go build hitbox.go

unset GOOS
unset GOARCH
```

## Configuration
Hitbox reads from a required `config.json` file at runtime.  The required fields are described below:

- `site_name`: The name of the website as you'd like to see it in your visitor notifications. A visitor's path (including leading `/`) is appended to this string.  Example value: `"site_name": "knik.co"`
- `slack_webhook_url`: The [incoming webhook URL](https://api.slack.com/incoming-webhooks) generated by Slack.  Follow Slack's instructions to configure an incoming webhook integration.  Example value: `"slack_webhook_url": "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX"`
- `geoip_url_format`: A `fmt.Sprintf()`-able format string denoting the URL to request GeoIP information from via GET request.  By default, Hitbox uses [freegeoip.net](http://freegeoip.net). Currently, GeoIP return values must be in JSON format with `city` (optional), `region` (optional), and `country` (required) fields.  Example value: `"geoip_url_format": "http://freegeoip.com/json/%s"`
- `ignore_referrer_domains`: If this field is left blank, all of your site's internal links will register as being reffered from the same domain. To reduce clutter, it's recommended that you include your site's domain in this array, along with any other referral domains that shouldn't be shown. Example value: `"ignore_referrer_domains": ["knik.co", "facebook.com"]`
- `port`: The port on which Hitbox should listen for GET requests. Example value: `"port": ":8080"`