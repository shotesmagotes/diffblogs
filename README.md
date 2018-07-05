# diffblogs
A backend augmenting a filesystem made for versioning and creating timeline driven news articles on web applications.

# How it works
Whenever a user adds/edits an article in a specified directory, a diff is taken between that article and any past articles of the same name.
The diff is used to construct a article text versioning system, which also compiles the history of the diffs into one article.
The compiled article allows the developer to showcase a timeline of revisions to the article.
