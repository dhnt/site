---
name: dhnt-site
description: Read the dhnt company pages (https://ai.dhnt.io/ and https://ai.dhnt.io/about) as structured content instead of scraping the browser app. Use when an agent needs what dhnt is, how to start with bashy, the team or the contact links, or must cite a specific section at a specific revision. Covers the Accept negotiation, the .md and .json twins, the ETag/revision header, declared section ids, and the citation record {url, section_id, revision, section_sha256, quote}.
compatibility: any HTTP client; `bashy fetch` or `curl` recommended
---

# dhnt-site — reading the dhnt pages

The pages are generated from the public repository
[github.com/dhnt/site](https://github.com/dhnt/site). Every page has three
representations of one content document: HTML (browsers), Markdown and JSON
(agents). They list the same links and the same section ids, and carry the
same revision.

## Pages

| Page | HTML | Markdown | JSON | Section ids |
|---|---|---|---|---|
| Home | `/` | `/index.md` | `/index.json` | `what`, `start`, `shape` |
| About | `/about` | `/about.md` | `/about.json` | `about`, `team`, `contact` |

Navigation, shared by both pages: Home `/`, About `/about`, GitHub
`https://github.com/dhnt`. Also served: `/robots.txt`, `/llms.txt`.

## How to read a page

1. Ask for the structure you want, either by `Accept` header or by twin URL:

   ```sh
   curl -sH 'Accept: text/markdown'    https://ai.dhnt.io/about
   curl -sH 'Accept: application/json' https://ai.dhnt.io/about
   curl -s https://ai.dhnt.io/about.json        # same as the line above
   bashy fetch --md https://ai.dhnt.io/about
   ```

   Anything else (including `*/*`) returns HTML. Every response carries
   `Vary: Accept`, a strong `ETag` equal to the page revision, and
   `Link: </about.md>; rel="alternate"; type="text/markdown"` (and the JSON
   twin likewise). HTML has `<link rel="alternate">` for both twins,
   `<link rel="canonical">` and `<meta name="revision">`.

2. JSON shape:

   ```json
   {"path": "/about", "title": "DHNT Inc.", "revision": "<rev>",
    "nav": [{"label": "Home", "href": "/"}],
    "sections": [{"id": "team", "heading": "Founder & team",
                  "html": "...", "markdown": "...", "sha256": "..."}]}
   ```

   Markdown starts with front matter (`path`, `revision`, `alternates`), then
   the navigation as labeled links, then one `## Heading {#id}` per section.

3. Section ids are declared in the source, never derived from headings: a
   heading may change while its id stays. Address a section as
   `https://ai.dhnt.io/about#team`.

## How to cite

Cite a section with the record

```json
{"url": "https://ai.dhnt.io/about#team", "section_id": "team",
 "revision": "<page revision>", "section_sha256": "<sections[].sha256>",
 "quote": "<exact text you rely on>"}
```

and in prose as `https://ai.dhnt.io/about#team @ <revision>`. The revision
changes when the copy changes, not when the rendering changes; a later
reader compares `section_sha256` to see whether the cited section moved.

## Shaping

"Shape this page" (`/?shape=open`) changes the page's look for a visitor; it
does not change the content. Content changes go through the owner's
`bashy shape` loop on this repository, citing the section and revision:
`bashy shape wish site --page /about --element 'section=team revision=<rev>' "<copy edit>"`.
A wish at a stale revision is refused ("stale, re-read").
