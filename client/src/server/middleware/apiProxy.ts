import { useRuntimeConfig } from "#imports";
import { defineEventHandler, H3Event, proxyRequest } from "h3";

const apiUrl = useRuntimeConfig().GO_API;
if (!apiUrl) {
  console.error("GO_API is required");
}

export default defineEventHandler((event: H3Event) => {
  if (event.node.req.url?.startsWith("/api/")) {
    return proxyRequest(event, `${apiUrl}${event.node.req.url}`);
  }
});
