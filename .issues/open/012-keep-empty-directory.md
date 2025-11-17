---
id: "012"
assignee: ""
labels: []
created: 2025-11-17T11:25:52.132547+09:00
updated: 2025-11-17T11:29:37.044212+09:00
---

# Keep empty directory

To make sure to keep empty open/closed directory in the git repo, we should add some hidden file. many project add .keep file for that.

## Todo

- [ ] add `.keep` files when we initialize directory like `.issues/open/.keep`
- [ ] Make sure open/close the last item in the directory won't delete directory as well as .keep file
