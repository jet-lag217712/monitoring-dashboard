export const paths = {
  home: () => '/',
  site: siteId => `/sites/${encodeURIComponent(siteId)}`,
  device: (siteId, deviceKey) =>
    `/sites/${encodeURIComponent(siteId)}/devices/${encodeURIComponent(deviceKey)}`,
  wall: () => '/wall',
  wallDisplay: () => '/wall/display',
}

export function isWallPath(pathname) {
  return pathname === paths.wall() || pathname === paths.wallDisplay()
}
