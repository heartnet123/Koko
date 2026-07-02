export interface Anime {
  mal_id: number
  title: string
  images: { jpg: { image_url: string; large_image_url: string } }
  score: number | null
  type: string | null
}
