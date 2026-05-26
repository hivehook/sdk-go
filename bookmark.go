package hivehook

import "context"

const bookmarkFragment = `
	id eventId name notes createdAt
`

const listBookmarksQuery = `query($eventId: UUID, $search: String, $limit: Int, $offset: Int, $after: String, $first: Int) {
	bookmarks(eventId: $eventId, search: $search, limit: $limit, offset: $offset, after: $after, first: $first) {
		nodes {` + bookmarkFragment + `}
		pageInfo { total limit offset endCursor hasNextPage }
	}
}`

const getBookmarkQuery = `query($id: UUID!) {
	bookmark(id: $id) {` + bookmarkFragment + `}
}`

const createBookmarkMutation = `mutation($eventId: UUID!, $name: String, $notes: String) {
	createBookmark(eventId: $eventId, name: $name, notes: $notes) {` + bookmarkFragment + `}
}`

const deleteBookmarkMutation = `mutation($id: UUID!) {
	deleteBookmark(id: $id)
}`

// BookmarkService manages saved Event Bookmarks.
type BookmarkService struct {
	gql *graphqlClient
}

// ListBookmarksOptions filters list results.
type ListBookmarksOptions struct {
	ListOptions
	EventID string
}

// List returns Bookmarks matching opts and a PageInfo cursor.
func (s *BookmarkService) List(ctx context.Context, opts *ListBookmarksOptions) ([]*Bookmark, *PageInfo, error) {
	vars := map[string]any{}
	if opts != nil {
		vars = opts.toVars()
		if opts.EventID != "" {
			vars["eventId"] = opts.EventID
		}
	}
	var result struct {
		Bookmarks struct {
			Nodes    []*Bookmark `json:"nodes"`
			PageInfo PageInfo    `json:"pageInfo"`
		} `json:"bookmarks"`
	}
	if err := s.gql.do(ctx, listBookmarksQuery, vars, &result); err != nil {
		return nil, nil, err
	}
	return result.Bookmarks.Nodes, &result.Bookmarks.PageInfo, nil
}

// Get fetches a single Bookmark by ID.
func (s *BookmarkService) Get(ctx context.Context, id string) (*Bookmark, error) {
	var result struct {
		Bookmark *Bookmark `json:"bookmark"`
	}
	if err := s.gql.do(ctx, getBookmarkQuery, map[string]any{"id": id}, &result); err != nil {
		return nil, err
	}
	return result.Bookmark, nil
}

// Create persists a new Bookmark.
func (s *BookmarkService) Create(ctx context.Context, eventID string, name, notes *string) (*Bookmark, error) {
	vars := map[string]any{
		"eventId": eventID,
	}
	if name != nil {
		vars["name"] = *name
	}
	if notes != nil {
		vars["notes"] = *notes
	}
	var result struct {
		CreateBookmark Bookmark `json:"createBookmark"`
	}
	if err := s.gql.do(ctx, createBookmarkMutation, vars, &result); err != nil {
		return nil, err
	}
	return &result.CreateBookmark, nil
}

// Delete removes a Bookmark by ID.
func (s *BookmarkService) Delete(ctx context.Context, id string) error {
	return s.gql.do(ctx, deleteBookmarkMutation, map[string]any{"id": id}, nil)
}
