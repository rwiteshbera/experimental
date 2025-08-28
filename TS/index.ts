import { createClient } from "@clickhouse/client";

const ch = createClient({
  url: "http://localhost:8123",
  username: "kloudmate",
  password: "hxx6HVlLoKvAAhz",
  keep_alive: { enabled: true },
});

export const client = {
  query: async (...args: Parameters<typeof ch.query>) => {
    try {
      return await ch.query(...args);
    } catch (error) {
      throw new Error("Query execution failed");
    }
  },
};

client
  .query({ query: "SHOW DATABASES", format: "JSONEachRow" })
  .then((data: any) => {
    console.log(data.query_id)
  })
  .catch((e) => {
    console.log(e.message);
  });
