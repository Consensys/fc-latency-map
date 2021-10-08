import type { NextApiResponse } from "next";
import type { AppProps, NextWebVitalsMetric } from "next/app";

export default async ({}, res: NextApiResponse) => {
  const formatMemoryUsage = (data: any) =>
    `${Math.round((data / 1024 / 1024) * 100) / 100} MB`;

  const memoryData = process.memoryUsage();

  const memoryUsage = {
    rssMemory: `${formatMemoryUsage(memoryData.rss)}`, // -> Resident Set Size - total memory allocated for the process execution
    heapMemoryTotal: `${formatMemoryUsage(memoryData.heapTotal)}`, // -> total size of the allocated heap
    heapMemoryUsed: `${formatMemoryUsage(memoryData.heapUsed)}`, // -> actual memory used during the execution
    externalMemory: `${formatMemoryUsage(memoryData.external)}`, // -> V8 external memory
  };

  res.status(200).json(memoryUsage);
};

export function reportWebVitals(metric: NextWebVitalsMetric) {
  console.log("XOXO::::", metric);
}
