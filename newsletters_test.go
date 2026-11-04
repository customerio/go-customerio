package customerio_test

import (
	"context"
	"testing"

	"github.com/customerio/go-customerio/v3"
)

func TestListNewsletters(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"newsletters": []map[string]any{{"id": 1, "name": "Weekly"}}})
	got, err := api.ListNewsletters(context.Background(), customerio.ListNewslettersOptions{Sort: customerio.SortDescending})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/newsletters?sort=desc")
	if len(got.Newsletters) != 1 || got.Newsletters[0].Name != "Weekly" {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestCreateNewsletter(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"newsletter": map[string]any{"id": 1}})
	data := map[string]any{"name": "Weekly"}
	got, err := api.CreateNewsletter(context.Background(), data)
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/newsletters")
	assertJSONEqual(t, rec.body, data)
	if got.Newsletter.ID != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestGetNewsletter(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetNewsletter(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "newsletterID")
	}

	api, rec := apiServer(t, 200, map[string]any{"newsletter": map[string]any{"id": 1}})
	if _, err := api.GetNewsletter(context.Background(), "1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/newsletters/1")
}

func TestDeleteNewsletter(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if err := api.DeleteNewsletter(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "newsletterID")
	}

	api, rec := apiServer(t, 204, nil)
	if err := api.DeleteNewsletter(context.Background(), "1"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "DELETE", "/v1/newsletters/1")
}

func TestNewsletterContents(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetNewsletterContents(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "newsletterID")
	}

	api, rec := apiServer(t, 200, map[string]any{"contents": []map[string]any{{"id": 5, "language": "en"}}})
	got, err := api.GetNewsletterContents(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/newsletters/1/contents")
	if len(got.Contents) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestNewsletterContent(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetNewsletterContent(context.Background(), "", "5"); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "newsletterID")
	}
	if _, err := api.GetNewsletterContent(context.Background(), "1", ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "contentID")
	}

	api, rec := apiServer(t, 200, map[string]any{"content": map[string]any{"id": 5}})
	if _, err := api.GetNewsletterContent(context.Background(), "1", "5"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/newsletters/1/contents/5")

	data := map[string]any{"subject": "hi"}
	if _, err := api.UpdateNewsletterContent(context.Background(), "1", "5", data); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/newsletters/1/contents/5")
	assertJSONEqual(t, rec.body, data)
}

func TestNewsletterContentMetrics(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"metric": map[string]any{"series": map[string]any{}}})
	if _, err := api.GetNewsletterContentMetrics(context.Background(), "1", "5", customerio.TransactionalMetricsOptions{Period: customerio.MetricsPeriodMonths, Steps: 1}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/newsletters/1/contents/5/metrics?period=months&steps=1")

	if _, err := api.GetNewsletterContentMetrics(context.Background(), "", "5", customerio.TransactionalMetricsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "newsletterID")
	}
	if _, err := api.GetNewsletterContentMetrics(context.Background(), "1", "", customerio.TransactionalMetricsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "contentID")
	}
}

func TestNewsletterContentMetricsLinks(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"links": []any{}})
	if _, err := api.GetNewsletterContentMetricsLinks(context.Background(), "1", "5", customerio.TransactionalLinkMetricsOptions{}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/newsletters/1/contents/5/metrics/links")
}

func TestNewsletterMetrics(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetNewsletterMetrics(context.Background(), "", customerio.TransactionalMetricsOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "newsletterID")
	}

	api, rec := apiServer(t, 200, map[string]any{"metric": map[string]any{"series": map[string]any{}}})
	if _, err := api.GetNewsletterMetrics(context.Background(), "1", customerio.TransactionalMetricsOptions{}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/newsletters/1/metrics")
}

func TestNewsletterMetricsLinks(t *testing.T) {
	api, rec := apiServer(t, 200, map[string]any{"links": []any{}})
	if _, err := api.GetNewsletterMetricsLinks(context.Background(), "1", customerio.TransactionalLinkMetricsOptions{}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/newsletters/1/metrics/links")
}

func TestGetNewsletterMessages(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetNewsletterMessages(context.Background(), "", customerio.NewsletterMessagesOptions{}); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "newsletterID")
	}

	api, rec := apiServer(t, 200, map[string]any{"messages": []map[string]any{{"id": "m1"}}})
	got, err := api.GetNewsletterMessages(context.Background(), "1", customerio.NewsletterMessagesOptions{Type: customerio.NewsletterChannelInbox})
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/newsletters/1/messages?type=inbox")
	if len(got.Messages) != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestSendNewsletter(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.SendNewsletter(context.Background(), "", nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "newsletterID")
	}

	api, rec := apiServer(t, 200, map[string]any{"newsletter": map[string]any{"id": 1}})
	if _, err := api.SendNewsletter(context.Background(), "1", map[string]any{}); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/newsletters/1/send")
}

func TestScheduleNewsletter(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.ScheduleNewsletter(context.Background(), "", nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "newsletterID")
	}

	api, rec := apiServer(t, 200, map[string]any{"newsletter": map[string]any{"id": 1}})
	data := map[string]any{"send_at": 100}
	if _, err := api.ScheduleNewsletter(context.Background(), "1", data); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/newsletters/1/schedule")
	assertJSONEqual(t, rec.body, data)
}

func TestNewsletterLanguage(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.CreateNewsletterLanguage(context.Background(), "", nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "newsletterID")
	}

	api, rec := apiServer(t, 200, map[string]any{"content": map[string]any{"language": "en"}})
	data := map[string]any{"language": "en"}
	if _, err := api.CreateNewsletterLanguage(context.Background(), "1", data); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/newsletters/1/language")
	assertJSONEqual(t, rec.body, data)

	if _, err := api.GetNewsletterLanguage(context.Background(), "1", "en"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/newsletters/1/language/en")

	if _, err := api.UpdateNewsletterLanguage(context.Background(), "1", "en", data); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/newsletters/1/language/en")

	if err := api.DeleteNewsletterLanguage(context.Background(), "1", "en"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "DELETE", "/v1/newsletters/1/language/en")

	if _, err := api.GetNewsletterLanguage(context.Background(), "1", ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "language")
	}
}

func TestNewsletterTestGroups(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.GetNewsletterTestGroups(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "newsletterID")
	}

	api, rec := apiServer(t, 200, map[string]any{"test_groups": []map[string]any{
		{"id": 1, "name": "A", "content_ids": []int{25, 26}},
	}})
	got, err := api.GetNewsletterTestGroups(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/newsletters/1/test_groups")
	if len(got.TestGroups) != 1 || len(got.TestGroups[0].ContentIDs) != 2 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestCreateNewsletterTestGroup(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.CreateNewsletterTestGroup(context.Background(), ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "newsletterID")
	}

	api, rec := apiServer(t, 200, map[string]any{"test_group": map[string]any{"id": 1, "name": "A"}})
	got, err := api.CreateNewsletterTestGroup(context.Background(), "1")
	if err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "POST", "/v1/newsletters/1/test_groups")
	assertJSONEqual(t, rec.body, nil)
	if got.TestGroup.ID != 1 {
		t.Errorf("unexpected response: %#v", got)
	}
}

func TestNewsletterTestGroupLanguage(t *testing.T) {
	api := customerio.NewAPIClient("myKey")
	if _, err := api.CreateNewsletterTestGroupLanguage(context.Background(), "", "1", nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "newsletterID")
	}
	if _, err := api.CreateNewsletterTestGroupLanguage(context.Background(), "1", "", nil); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "testGroupID")
	}

	api, rec := apiServer(t, 200, map[string]any{"content": map[string]any{"language": "en"}})
	data := map[string]any{"language": "en"}
	if _, err := api.CreateNewsletterTestGroupLanguage(context.Background(), "1", "1", data); err != nil {
		t.Fatal(err)
	}
	// Wire path uses singular "test_group", unlike the plural "test_groups" list/create endpoints.
	assertAPIRequest(t, rec, "POST", "/v1/newsletters/1/test_group/1/language")
	assertJSONEqual(t, rec.body, data)

	if _, err := api.GetNewsletterTestGroupLanguage(context.Background(), "1", "1", "en"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "GET", "/v1/newsletters/1/test_group/1/language/en")

	if _, err := api.UpdateNewsletterTestGroupLanguage(context.Background(), "1", "1", "en", data); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "PUT", "/v1/newsletters/1/test_group/1/language/en")

	if err := api.DeleteNewsletterTestGroupLanguage(context.Background(), "1", "1", "en"); err != nil {
		t.Fatal(err)
	}
	assertAPIRequest(t, rec, "DELETE", "/v1/newsletters/1/test_group/1/language/en")

	if _, err := api.GetNewsletterTestGroupLanguage(context.Background(), "1", "1", ""); err == nil {
		t.Fatal("expected error")
	} else {
		checkParamError(t, err, "language")
	}
}
