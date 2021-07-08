# Customerio

# go-customerio [![CircleCI](https://circleci.com/gh/customerio/go-customerio/tree/master.svg?style=svg)](https://circleci.com/gh/customerio/go-customerio/tree/master)

A golang client for the [Customer.io](http://customer.io) [event API](https://app.customer.io/api/docs/index.html).
_Tested with Go1.16_

Godoc here: [https://godoc.org/github.com/customerio/go-customerio](https://godoc.org/github.com/customerio/go-customerio)

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

    $ go get github.com/customerio/go-customerio

## Usage

### Before we get started: API client vs. JavaScript snippet

It's helpful to know that everything below can also be accomplished
through the [Customer.io JavaScript snippet](http://customer.io/docs/basic-integration.html).

In many cases, using the JavaScript snippet will be easier to integrate with
your app, but there are several reasons why using the API client is useful:

- You're not planning on triggering emails based on how customers interact with
  your website (e.g. users who haven't visited the site in X days)
- You're using the javascript snippet, but have a few events you'd like to
  send from your backend system. They will work well together!
- You'd rather not have another javascript snippet slowing down your frontend.
  Our snippet is asynchronous (doesn't affect initial page load) and very small, but we understand.

In the end, the decision on whether or not to use the API client or
the JavaScript snippet should be based on what works best for you.
You'll be able to integrate **fully** with [Customer.io](http://customer.io) with either approach.

### Setup

Create an instance of the client with your [Customer.io credentials](https://fly.customer.io/settings/api_credentials).

```go
track := customerio.NewTrackClient("YOUR SITE ID", "YOUR API SECRET KEY", customerio.WithRegion(customerio.RegionUS))
```

Your account region—`RegionUS` or `RegionEU`—is optional. If you do not specify your region, we assume that your account is based in the US (`RegionUS`). If your account is based in the EU and you do not provide the correct region, we'll route requests from the US to `RegionEU` accordingly, however this may cause data to be logged in the US. 

### Identify logged in customers

Tracking data of logged in customers is a key part of [Customer.io](http://customer.io). In order to
send triggered emails, we must know the email address of the customer. You can
also specify any number of customer attributes which help tailor [Customer.io](http://customer.io) to your
business.

Attributes you specify are useful in several ways:

- As customer variables in your triggered emails. For instance, if you specify
  the customer's name, you can personalize the triggered email by using it in the
  subject or body.

- As a way to filter who should receive a triggered email. For instance,
  if you pass along the current subscription plan (free / basic / premium) for your customers, you can
  set up triggers which are only sent to customers who have subscribed to a
  particular plan (e.g. "premium").

You'll want to identify your customers when they sign up for your app and any time their
key information changes. This keeps [Customer.io](http://customer.io) up to date with your customer information.

```go
// Arguments
// customerID (required) - a unique identifier string for this customers
// attributes (required) - a ```map[string]interface{}``` of information about the customer. You can pass any
//                         information that would be useful in your triggers. You
//                         should at least pass in an email, and created_at timestamp.
//                         your interface{} should be parseable as Json by 'encoding/json'.Marshal

track.Identify("5", map[string]interface{}{
  "email": "bob@example.com",
  "created_at": time.Now().Unix(),
  "first_name": "Bob",
  "plan": "basic",
})
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

track.Delete("5")
```

### Tracking a custom event

Now that you're identifying your customers with [Customer.io](http://customer.io), you can now send events like
"purchased" or "watchedIntroVideo". These allow you to more specifically target your users
with automated emails, and track conversions when you're sending automated emails to
encourage your customers to perform an action.

```go
// Arguments
// customerID (required)  - the id of the customer who you want to associate with the event.
// name (required)        - the name of the event you want to track.
// attributes (optional)  - any related information you'd like to attach to this
//                          event, as a ```map[string]interface{}```. These attributes can be used in your triggers to control who should
//                         receive the triggered email. You can set any number of data values.

track.Track("5", "purchase", map[string]interface{}{
    "type": "socks",
    "price": "13.99",
})
```

### Tracking an anonymous event

You can also send anonymous events representing people you haven't identified. An anonymous event requires an `anonymous_id` representing the unknown person and an event `name`. When you identify a person, you can set their `anonymous_id` attribute. If [event merging](https://customer.io/docs/anonymous-events/#turn-on-merging) is turned on in your workspace, and the attribute matches the `anonymous_id` in one or more events that were logged within the last 30 days, we associate those events with the person.

```go
// Arguments
// anonymous_id (required)    - an identifier representing an unknown person.
// name (required)            - the name of the event you want to track.
// attributes (optional)      - any related information you'd like to attach to this
//                              event, as a ```map[string]interface{}```. These attributes can be used in your triggers to control who should
//                              receive the triggered email. You can set any number of data values.

track.TrackAnonymous("anonymous_id", "invite", map[string]interface{}{
    "first_name": "Alex",
    "source": "OldApp",
})
```

### Adding a device to a customer

In order to send push notifications, we need customer device information.

```go
// Arguments
// customerID (required) - a unique identifier string for this customer
// deviceID (required)   - a unique identifier string for this device
// platform (required)   - the platform of the device, currently only accepts 'ios' and 'andriod'
// data (optional)        - a ```map[string]interface{}``` of information about the device. You can pass any
//                         key/value pairs that would be useful in your triggers. We
//                         currently only save 'last_used'.
//                         your interface{} should be parseable as Json by 'encoding/json'.Marshal

track.AddDevice("5", "messaging token", "android", map[string]interface{}{
"last_used": time.Now().Unix(),
})
```

### Deleting devices

Deleting a device will remove it from the customers device list in Customer.io.

```go
// Arguments
// customerID (required) - the id of the customer the device you want to delete belongs to
// deviceToken (required) - a unique identifier for the device.  This
//                          should be the same id you'd pass into the
//                          `addDevice` command above

track.DeleteDevice("5", "messaging-token")
```

### Send Transactional Messages

To use the Customer.io [Transactional API](https://customer.io/docs/transactional-api), create an instance of the API client using an [app key](https://customer.io/docs/managing-credentials#app-api-keys).

Create a `SendEmailRequest` instance, and then use `SendEmail` to send your message. [Learn more about transactional messages and optional `SendEmailRequest` properties](https://customer.io/docs/transactional-api).

You can also send attachments with your message. Use `Attach` to encode attachments.

```go
import "github.com/customerio/go-customerio"

client := customerio.NewAPIClient("<extapikey>", customerio.WithRegion(customerio.RegionUS));

// TransactionalMessageId — the ID of the transactional message you want to send.
// To                     — the email address of your recipients.
// Identifiers            — contains the id of your recipient. If the id does not exist, Customer.io creates it.
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
  fmt.Println(err)
}
request.Attach("receipt.pdf", f)

body, err := client.SendEmail(context.Background(), &request)
if err != nil {
  fmt.Println(err)
}

fmt.Println(body)
```

## Contributing

1. Fork it
2. Clone your fork (`git clone git@github.com:MY_USERNAME/go-customerio.git && cd go-customerio`)
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Commit your changes (`git commit -am 'Added some feature'`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create new Pull Request
