# Customer.io Segments API Methods

This document explains how to use the new methods added for managing customer segments in Customer.io. These methods enable you to list, create, delete segments, and manage customers in segments.

## Usage

### Available Methods

1. **Create Segment**

   Creates a new segment in Customer.io.

    ```go
    request := &CreateSegmentRequest{
		Segment: Segment{
			Name: "New Segment",
		},
	}

	resp, err := cio.CreateSegment(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Segment created: ID: %d, Name: %s", resp.Segment.ID, resp.Segment.Name)
    ```

2. **List Segments**

   Retrieves a list of all customer segments.

    ```go
    segments, err := cio.ListSegments(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    for _, segment := range segments.Segments {
        log.Printf("Segment ID: %d, Name: %s", segment.ID, segment.Name)
    }
    ```

3. **Get Segment by ID**

   Retrieves a specific segment by its ID.

    ```go
    segmentID := 1234
	resp, err := cio.GetSegment(context.Background(), segmentID)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Segment ID: %d, Name: %s", resp.Segment.ID, resp.Segment.Name)
    ```

4. **Delete Segment**

   Deletes a segment by its ID.

    ```go
    segmentID := 1234
    err := cio.DeleteSegment(context.Background(), segmentID)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Segment deleted successfully")
    ```

5. **Get Segment Dependencies**

   Retrieves dependencies for a specific segment.

    ```go
    segmentID := 1234
    dependencies, err := cio.GetSegmentDependencies(context.Background(), segmentID)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Segment dependencies: %v", dependencies)
    ```

6. **Get Segment Customer Count**

   Retrieves the number of customers in a specific segment.

    ```go
    segmentID := 1234
    count, err := cio.GetSegmentCustomerCount(context.Background(), segmentID)
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Customer count in segment %d: %d", segmentID, count.Count)
    ```

7. **List Customers in a Segment**

   Retrieves a list of customers in a specific segment.

    ```go
    segmentID := 1234
	customers, err := cio.ListCustomersInSegment(context.Background(), segmentID)
	if err != nil {
		log.Fatal(err)
	}

	for _, customer := range customers.Identifiers {
		log.Printf("Customer ID: %d", customer.ID)
	}
    ```

### Managing Customers in Segments

You can add or remove customers from segments using the following methods:

1. **Add People to Segment**

   Adds a list of customer IDs to a segment.

    ```go
    segmentID := 1234
	customerIDs := []string{"customer_1", "customer_2"}

	err := track.AddPeopleToSegment(context.Background(), segmentID, customerIDs)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Customers added to segment successfully")
    ```

2. **Remove People from Segment**

   Removes a list of customer IDs from a segment.

    ```go
    segmentID := 1234
	customerIDs := []string{"customer_1", "customer_2"}

	err := track.RemovePeopleFromSegment(context.Background(), segmentID, customerIDs)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Customers removed from segment successfully")
    ```

## Example: Creating a Segment and Adding Customers

```go
func main() {
    // Create a new segment
	request := &CreateSegmentRequest{
		Segment: Segment{
			Name: "VIP Customers",
		},
	}
	resp, err := cio.CreateSegment(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Created segment: %s (ID: %d)", resp.Segment.Name, resp.Segment.ID)

	// Add customers to the new segment
	customerIDs := []string{"customer_1", "customer_2"}
	err = track.AddPeopleToSegment(context.Background(), resp.Segment.ID, customerIDs)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Customers added to the segment successfully")
}
