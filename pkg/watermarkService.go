package watermark

import (
	"context"
	"net/http"
	"os"
	"watermark/internal"

	"github.com/go-kit/log"
	"github.com/lithammer/shortuuid"
)

type watermarkService struct{}

func NewService() Service {
	return &watermarkService{}
}

func (w *watermarkService) Get(_ context.Context, filters ...internal.Filter) ([]internal.Document, error) {
	// Query the database using the filters and return the list of documents.
	// Return an error if the filter (key) is invalid or if no item is found.
	doc := internal.Document{
		Content: "book",
		Title:   "Harry Potter and the Half-Blood Prince",
		Author:  "J.K. Rowling",
		Topic:   "Fiction and Magic",
	}
	return []internal.Document{doc}, nil
}

func (w *watermarkService) Status(_ context.Context, ticketID string) (internal.Status, error) {
	// Query the database using the ticketID and return the document info.
	// Return an error if the ticketID is invalid or no document exists for that ticketID.
	return internal.InProgress, nil
}

func (w *watermarkService) Watermark(_ context.Context, ticketID, mark string) (int, error) {
	// Update the database entry with the watermark field as non-empty.
	// First, check if the watermark status is not already in InProgress, Started, or Finished state.
	// If yes, then return an invalid request.
	// Return an error if no item is found using the ticketID.
	return http.StatusOK, nil
}

func (w *watermarkService) AddDocument(_ context.Context, doc *internal.Document) (string, error) {
	// Add the document entry in the database by calling the database service.
	// Return an error if the doc is invalid and/or if there's a database entry error.
	newTicketID := shortuuid.New()
	return newTicketID, nil
}

func (w *watermarkService) ServiceStatus(_ context.Context) (int, error) {
	logger.Log("Checking the Service health...")
	return http.StatusOK, nil
}

var logger log.Logger

func init() {
	logger = log.NewLogfmtLogger(log.NewSyncWriter(os.Stderr))
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
}
