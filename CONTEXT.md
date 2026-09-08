# Koko

Anime discovery, personal tracking, and episode streaming web application.

## Language

**Anime**:
An animated production catalog item uniquely identified by its MyAnimeList identifier (`mal_id`).
_Avoid_: Show, Series, Title, Media

**User**:
A registered account holder with credentials and profile details.
_Avoid_: Account, Profile, Member, Customer

**Watchlist**:
A user's personal collection of bookmarked anime saved for future reference. Binary membership only.
_Avoid_: Queue, Favorites, Tracker, Playlist

**Episode**:
A single numbered installment of an **Anime** available for video playback.
_Avoid_: Chapter, Part, Clip, Video

**Stream Source**:
A provider endpoint or embed delivering video playback for an **Episode**.
_Avoid_: Mirror, Server, Feed, Link

**Watch Progress**:
The recorded playback timestamp and completion status for a specific **Episode**.
_Avoid_: Watch History, Playback State, Resume Point

## Relationships

- An anonymous guest can browse, search, and stream without being a **User**
- A **User** owns exactly one **Watchlist**
- A **Watchlist** contains zero or more **Anime**
- An **Anime** contains one or more **Episodes**
- An **Episode** resolves to one or more **Stream Sources**
- A **User** records **Watch Progress** per **Episode** (persisted in DB; guests persist locally)
- An **Anime** can appear in many **Watchlists**

## Example dialogue

> **Dev:** "Does **Watchlist** track which episode the user is on?"
> **Domain expert:** "No. **Watchlist** is strictly bookmarking. Playback position is tracked separately by **Watch Progress**."
> **Dev:** "Where does **Watch Progress** live for guests?"
> **Domain expert:** "Guests keep **Watch Progress** in local browser storage; authenticated **Users** sync it to the database."

## Flagged ambiguities

- "Show" or "Series" used interchangeably with **Anime** — resolved: **Anime** is canonical term; records are external and read-only.
- **Watchlist** confused with multi-state progress tracking — resolved: **Watchlist** is binary bookmark list only.
- "Account" vs "Profile" vs **User** — resolved: unified under **User**. No separate profile entity.
- "Server" vs "Mirror" vs "Stream" — resolved: **Stream Source** is canonical term.
- "History" vs **Watch Progress** — resolved: **Watch Progress** is canonical; captures timestamp and completion.
