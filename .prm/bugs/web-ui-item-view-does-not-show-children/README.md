# Web UI item view does not show children

## Problem
The entity detail page in the web UI only displays a count of children (e.g. '3 items') but does not render the actual child entities. Users cannot see or navigate to children from the item view.

## Expected Behavior
When viewing an epic or story that has children, the detail page should display a list of child entities with their type, title, status, and priority — each linking to the child's detail page.

## Fix
- Backend: resolve child slugs to full entities in the GET /api/entities/{id} handler
- Frontend: render a Children section with clickable rows between the metadata grid and the README