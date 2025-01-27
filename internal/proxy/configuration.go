package proxy

// ============================================
// Websocket configuration constants
// ============================================
const wsReadBufferSize = 1024
const wsWriteBufferSize = 1024

// limits the max size of a websocket packet to 10MB so someone can't just crash the server.
// The proxy clients break packets up in to 1MB chunks, so under normal operation this should never be hit.
const wsMaxPacketSize = 1024 * 1024 * 10 // 10MB

const proxyChannelBufferSize = 1024
