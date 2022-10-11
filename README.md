<p align="center">
  <a href="https://customer.io">
    <img src="https://user-images.githubusercontent.com/6409227/144680509-907ee093-d7ad-4a9c-b0a5-f640eeb060cd.png" height="60">
  </a>
  <p align="center">Power automated communication that people like to receive.</p>
</p>

![Latest release](https://img.shields.io/github/v/release/customerio/go-customerio)
![Software License](https://img.shields.io/github/license/customerio/go-customerio)
[![CI status](https://github.com/customerio/go-customerio/actions/workflows/main.yml/badge.svg)](https://github.com/customerio/go-customerio/actions/workflows/main.yml)
[![codecov](https://codecov.io/gh/customerio/go-customerio/branch/main/graph/badge.svg?token=D59CJnFVDV)](https://codecov.io/gh/customerio/go-customerio)
![Go version](https://img.shields.io/github/go-mod/go-version/customerio/go-customerio)
[![Go Doc](https://img.shields.io/badge/Go_Doc-reference-blue.svg)](https://pkg.go.dev/github.com/customerio/go-customerio/v3)

# Customer.io Go 

A Go client library for interacting with the [Customer.io API](https://customer.io/docs/api/).

## Installation

Add this line to your application's imports:

```go
import (
    // ...
    "github.com/customerio/go-customerio/v3"
)
```

And then execute:

    go get

Or install it yourself:

    $ go get github.com/customerio/go-customerio/v3

## Before we get started: API client vs. JavaScript snippet

It's helpful to know that everything below can also be accomplished
through the [Customer.io JavaScript snippet](https://customer.io/docs/basic-integration.html).

In many cases, using the JavaScript snippet will be easier to integrate with
your app, but there are several reasons why using the API client is useful:

- You're not planning on triggering emails based on how customers interact with
  your website (e.g. users who haven't visited the site in X days)
- You're using the JavaScript snippet, but have a few events you'd like to
  send from your backend system. They will work well together!
- You'd rather not have another JavaScript snippet slowing down your frontend.
  Our snippet is asynchronous (doesn't affect initial page load) and very small, but we understand.

In the end, the decision on whether or not to use the API client or
the JavaScript snippet should be based on what works best for you.
You'll be able to integrate **fully** with [Customer.io](https://customer.io) with either approach.

Create an instance of the Track API or App API client with your [Customer.io credentials](https://fly.customer.io/settings/api_credentials).

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/customerio/go-customerio/v3"
)

var (
	// You can find or create new API credentials in your Customer.io
	// account under "Account Settings" => "API Credentials":
	// (https://fly.customer.io/settings/api_credentials)
	siteID      string = "your-site-id"
	trackAPIKey string = "your-track-api-key"
	appAPIKey   string = "your-app-api-key"
)

func main() {
	// Create an instance of the Customer.io Track API client
	track := customerio.NewTrackClient(siteID, trackAPIKey, customerio.WithRegion(customerio.RegionUS))

	// Send an Identify TrackAPI call
	if err := track.Identify("5", map[string]interface{}{
		"email":      "lucy@example.com",
		"created_at": time.Now().Unix(),
		"first_name": "Lucy",
		"plan":       "basic",
	}); err != nil {
		log.Fatal(err)
	}

	// Create an instance of the Customer.io App API Client
	cio := customerio.NewAPIClient(appAPIKey, customerio.WithRegion(customerio.RegionUS))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// The request object allows you to specify recipients and message data
	request := customerio.SendEmailRequest{
		Identifiers: map[string]string{
			"id": "customer_1",
		},
		To:      "customer@example.com",
		From:    "business@example.com",
		Subject: "hello, {{ trigger.name }}",
		Body:    "hello from the Customer.io {{ trigger.client }} client",
		MessageData: map[string]interface{}{
			"client": "Go",
			"name":   "gopher",
		},
	}

	// Send the email with a 10 second timeout
	resp, err := cio.SendEmail(ctx, &request)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Resp: %s\n", resp)
}

```

Your account region (`customerio.RegionUS` or `customerio.RegionEU`) is optional. If you do not specify your region, we assume that your account is based in the US (`customerio.RegionUS`). 

If your account is based in the EU and you do not provide the correct region, we'll route requests from the US to `customerio.RegionEU` accordingly, however this may cause data to be logged in the US. 

### Identify logged in customers

Tracking data of logged in customers is a key part of [Customer.io](https://customer.io). In order to send triggered messages, we must know the email address of the customer to send email or the phone number for SMS.

You can also specify any number of customer attributes which help tailor [Customer.io](https://customer.io) to your business.

Attributes you specify are useful in several ways:

- As customer variables in your triggered messages. For instance, if you specify the customer's name, you can personalize the triggered message by using it in the subject or body.

- As a way to filter who should receive a triggered message. For instance, if you pass along the current subscription plan (free / basic / premium) for your customers, you can set up triggers which are only sent to customers who have subscribed to a particular plan (e.g. "premium").

You'll want to identify your customers when they sign up for your product and any time their key information changes. This keeps [Customer.io](https://customer.io) up to date with your customer information.

```go
// Arguments
// customerID (required) - a unique identifier string for this customers
// attributes (required) - a ```map[string]interface{}``` of information about the customer. You can pass any
//                         information that would be useful in your triggers. You
//                         should at least pass in an email, and created_at timestamp.
//                         your interface{} should be parseable as JSON by 'encoding/json'.Marshal

if err := track.Identify("5", map[string]interface{}{
  "email": "bob@example.com",
  "created_at": time.Now().Unix(),
  "first_name": "Bob",
  "plan": "basic",
}); err != nil {
  // handle error
}
```

### Deleting customers

Deleting a customer will remove them, and all their information from
Customer.io. Note: if you're still sending data to Customer.io via
other means (such as the javascript snippet), the customer could be
recreated.

```go
// Arguments
// customerID (required) - a unique identifier for the customer.  This
//                          should be the same id you'd pass into the
//                          `identify` command above.

if err := track.Delete("5"); err != nil {
  // handle error
}
```

### Merge Duplicate Customers

When you merge two people, you pick a primary person and merge a secondary, duplicate person into the primary person. The primary person remains after the merge and the secondary person is deleted. This process is permanent: you cannot recover the secondary person.

The first and third parameters represent the identifier for the primary and secondary people respectively, one of `id`, `email`, or `cio_id`. The second and fourth parameters are the identifier values for the primary and secondary people, respectively.

```go
if err := track.MergeCustomers(customerio.IdentifierTypeEmail, "cool.person@company.com", customerio.IdentifierTypeCioID, "C123"); err != nil {
  // handle error
}
```

### Tracking a custom event

Now that you're identifying your customers with [Customer.io](https://customer.io), you can now send events like "purchased" or "watchedIntroVideo". 

These allow you to more specifically target your users with automated messages, and track conversions when you're sending automated messages to encourage your customers to perform an action.

```go
// Arguments
// customerID (required)  - the id of the customer who you want to associate with the event.
// name (required)        - the name of the event you want to track.
// attributes (optional)  - any related information you'd like to attach to this
//                          event, as a ```map[string]interface{}```. 
//                          These attributes can be used in your triggers to control who should
//                          receive the triggered message. You can set any number of data values.

if err := track.Track("5", "purchase", map[string]interface{}{
    "type": "socks",
    "price": "13.99",
}); err != nil {
  // handle error
}
```

### Tracking an anonymous event

You can also send anonymous events representing people you haven't identified. An anonymous event requires an `anonymous_id` representing the unknown person and an event `name`. When you identify a person, you can set their `anonymous_id` attribute. If [event merging](https://customer.io/docs/anonymous-events/#turn-on-merging) is turned on in your workspace, and the attribute matches the `anonymous_id` in one or more events that were logged within the last 30 days, we associate those events with the person.

```go
// Arguments
// anonymous_id (required)    - nullable, an identifier representing an unknown person.
// name (required)            - the name of the event you want to track.
// attributes (optional)      - any related information you'd like to attach to this
//                              event, as a ```map[string]interface{}```. 
//                              These attributes can be used in your triggers to control who should
//                              receive the triggered message. You can set any number of data values.

if err := track.TrackAnonymous("anonymous_id", "invite", map[string]interface{}{
    "first_name": "Alex",
    "source": "OldApp",
}); err != nil {
  // handle error
}
```
#### Anonymous invite events

If you previously sent [invite events](https://customer.io/docs/anonymous-invite-emails/), you can achieve the same functionality by sending an anonymous event an empty string for the anonymous identifier. To send anonymous invites, your event *must* include a `recipient` attribute. 

```go
if err := track.TrackAnonymous("", "invite", map[string]interface{}{
    "first_name": "Alex",
    "recipient": "alex.person@example.com",
}); err != nil {
  // handle error
}
```

### Adding a device to a customer

In order to send push notifications, we need customer device information.

```go
// Arguments
// customerID (required) - a unique identifier string for this customer
// deviceID (required)   - a unique identifier string for this device
// platform (required)   - the platform of the device, currently only accepts 'ios' and 'android'
// data (optional)       - a ```map[string]interface{}``` of information about the device. 
//                         You can pass any key/value pairs that would be useful in your triggers. 
//                         We currently only save 'last_used'.
//                         Your interface{} should be parseable as Json by 'encoding/json'.Marshal

if err := track.AddDevice("5", "messaging token", "android", map[string]interface{}{
"last_used": time.Now().Unix(),
}); err != nil {
  // handle error
}
```

### Deleting devices

Deleting a device will remove it from the customer's device list in Customer.io.

```go
// Arguments
// customerID (required)  - the id of the customer the device you want to delete belongs to
// deviceToken (required) - a unique identifier for the device.
//                          This should be the same id you'd pass into the
//                          `addDevice` command above

if err := track.DeleteDevice("5", "messaging-token"); err != nil {
  // handle error
}
```

### Send Transactional Messages

To use the Customer.io [Transactional API](https://customer.io/docs/transactional-api), create an instance of the API client using an [App API key](https://customer.io/docs/managing-credentials#app-api-keys).

Create a `customerio.SendEmailRequest` instance, and then use `(c *customerio.APIClient).SendEmail` to send your message. [Learn more about transactional messages and optional `SendEmailRequest` properties](https://customer.io/docs/transactional-api).

You can also send attachments with your message. Use `customerio.SendEmailRequest.Attach` to encode attachments.

```go
client := customerio.NewAPIClient("<extapikey>", customerio.WithRegion(customerio.RegionUS));

// TransactionalMessageId — the ID of the transactional message you want to send.
// To                     — the email address of your recipients.
// Identifiers            — contains the id of your recipient. 
//                          If the id does not exist, Customer.io creates it.
// MessageData            — contains properties that you want reference in your message using liquid.
// Attach                 — a helper that encodes attachments to your message.

request := customerio.SendEmailRequest{
  To: "person@example.com",
  TransactionalMessageID: "3",
  MessageData: map[string]interface{}{
    "name": "Person",
    "items": map[string]interface{}{
      "name": "shoes",
      "price": "59.99",
    },
    "products": []interface{}{},
  },
  Identifiers: map[string]string{
    "id": "example1",
  },
}

// (optional) attach a file to your message.
f, err := os.Open("receipt.pdf")
if err != nil {
  // handle error
}
defer f.Close()

request.Attach("receipt.pdf", f)

body, err := client.SendEmail(context.Background(), &request)
if err != nil {
  // handle error
}

fmt.Println(body)
```

## Context Support
There are additional API methods that support passing a context that satisfies the `context.Context` interface to allow better control over dispatched requests. For example with sending an event:
```go
// Create an instance of the Customer.io Track API client
track := customerio.NewTrackClient(siteID, trackAPIKey, customerio.WithRegion(customerio.RegionUS))

// Create a context with a 5s deadline
ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
defer cancel()

if err := track.TrackCtx(ctx, "5", "purchase", map[string]interface{}{
    "type": "socks",
    "price": "13.99",
}); err != nil {
  // handle error
}
```

## Contributing

1. Fork it
2. Clone your fork (`git clone git@github.com:MY_USERNAME/go-customerio.git && cd go-customerio`)
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Commit your changes (`git commit -am 'feat: Added some feature'`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create new Pull Request
