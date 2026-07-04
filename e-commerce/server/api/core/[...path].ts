import { defineEventHandler, proxyRequest } from "h3";

export default defineEventHandler(async (event) => {
  const config = useRuntimeConfig();
  const targetBase =
    process.env.SERVICE_CORE_API_URL ||
    config.public.serviceCoreApiUrl ||
    "http://127.0.0.1:7129";
  const subPath = event.context.params?.path || "";

  const targetUrl = `${targetBase.replace(/\/$/, "")}/${subPath.replace(/^\//, "")}`;
  console.log(
    `[PROXY] Proxying request from ${event.node.req.url} to ${targetUrl}`,
  );

  try {
    return await proxyRequest(event, targetUrl);
  } catch (err: any) {
    console.error(`[PROXY ERROR] Failed to proxy to ${targetUrl}:`, err);
    throw err;
  }
});
