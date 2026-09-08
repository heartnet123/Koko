import test from 'node:test'
import assert from 'node:assert/strict'
import fs from 'node:fs'
import path from 'node:path'
import { execSync } from 'node:child_process'

test('Integration: Go backend episodes endpoint and caching', () => {
  const output = execSync('go test -v ./...', { cwd: './backend', encoding: 'utf-8' })
  assert.match(output, /PASS: TestEpisodesEndpointAndCache/, 'TestEpisodesEndpointAndCache must pass')
  assert.match(output, /PASS: TestCacheOperations/, 'TestCacheOperations must pass')
})

test('Integration: Details page renders Episodes list and links to /watch/:animeId/:episodeNumber', () => {
  const detailsPath = path.resolve('./frontend/app/pages/movie/[id].vue')
  assert.ok(fs.existsSync(detailsPath), 'movie/[id].vue must exist')
  const content = fs.readFileSync(detailsPath, 'utf-8')

  assert.match(content, /\/api\/anime\/\$\{animeId\.value\}\/episodes/, 'fetches cached episodes from backend proxy')
  assert.match(content, /id="episodes-section"/, 'contains episodes section')
  assert.match(content, /:to="`\/watch\/\$\{animeId\}\/\$\{ep\.mal_id\}`"/, 'links to /watch/:animeId/:episodeNumber route')
  assert.match(content, /formatEpisodeDate/, 'formats episode air date')
  assert.match(content, /episodesPagination/, 'supports pagination')
})

test('Integration: Watch page route /watch/:animeId/:episodeNumber exists and handles playback navigation', () => {
  const watchPath = path.resolve('./frontend/app/pages/watch/[animeId]/[episodeNumber].vue')
  assert.ok(fs.existsSync(watchPath), 'watch/[animeId]/[episodeNumber].vue must exist')
  const content = fs.readFileSync(watchPath, 'utf-8')

  assert.match(content, /route\.params\.animeId/, 'parses animeId param')
  assert.match(content, /route\.params\.episodeNumber/, 'parses episodeNumber param')
  assert.match(content, /:to="`\/movie\/\$\{animeId\}`"/, 'navigates back to anime details')
  assert.match(content, /Stream Source/, 'includes stream player surface')
  assert.match(content, /Next Ep/, 'provides episode switcher controls')
})

test('Integration: Nuxt build successfully outputs compiled chunks for details and watch pages', () => {
  const buildDir = path.resolve('./frontend/.output/server/chunks/build')
  assert.ok(fs.existsSync(buildDir), 'Nuxt build output directory must exist')
  const files = fs.readdirSync(buildDir)

  const hasEpisodeChunk = files.some(f => f.includes('_episodeNumber_'))
  const hasDetailsChunk = files.some(f => f.includes('_id_'))

  assert.ok(hasEpisodeChunk, 'Build output should contain compiled watch chunk (_episodeNumber_)')
  assert.ok(hasDetailsChunk, 'Build output should contain compiled details chunk (_id_)')
})
