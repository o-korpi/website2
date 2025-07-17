---
title: Todo
wpm: 300
draft: false
created: 2025-07-15
desc: Website TODO
---

## Footnotes
Using the following syntax
`{footnote: This is a footnote}`
Could be replaced by sidenotes?
Bonus with footnotes would probably be automatic numbering, suitable for references.


## Sidenotes
Using the following syntax
`{sidenote: This is a sidenote!}`
or {sidenote This is a sidenote, for example.}
```
{sidenote}
This is a sidenote block
{/sidenote}
```

Should use `<aside>` element.
Example to the right! {sidenote Hello world!}

And next to this paragraph there should be a big one.
{sidenote}
Hello world! Lorem ipsum dolor sit amet.
{/sidenote}

Test.

Lorem ipsum dolor sit amet. Abc. {sidenote Test sidenote 1. Maybe with some really long text too} Test with two sidenotes 
te same paragraph. Will this work? {sidenote Abc}

## Table of contents
Self-explanatory.
Parse headers

## Frontmatter data display

- Use frontmatter title
- Make an estimated read time from WPM
- Show created and updated
- Show descripton
- Improved author display
- Leave drafts for later

## Charts

Using the following syntax
`{chart data="mydata.json" type="line" width="600"}`

{chart}

## Caching

Self-explanatory

Could use hashes to detect changes
