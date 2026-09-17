import { describe, expect, it } from 'vitest'
import { isWallPath, paths } from './paths.js'

describe('paths', () => {
  it('builds home, wall, and display routes', () => {
    expect(paths.home()).toBe('/')
    expect(paths.wall()).toBe('/wall')
    expect(paths.wallDisplay()).toBe('/wall/display')
  })

  it('encodes site and device keys', () => {
    expect(paths.site('site-a')).toBe('/sites/site-a')
    expect(paths.device('site-a', '10.255.1.1')).toBe('/sites/site-a/devices/10.255.1.1')
    expect(paths.device('site a', 'core/1')).toBe('/sites/site%20a/devices/core%2F1')
  })

  it('detects wall editor and display paths', () => {
    expect(isWallPath('/wall')).toBe(true)
    expect(isWallPath('/wall/display')).toBe(true)
    expect(isWallPath('/')).toBe(false)
    expect(isWallPath('/sites/site-a')).toBe(false)
  })
})
