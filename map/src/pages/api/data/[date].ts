import type { NextApiRequest, NextApiResponse } from "next";
import * as fs from "fs";

export default async (req: NextApiRequest, res: NextApiResponse) => {
  const { date } = req.query;
  const filepath = `/home/rapha/studio/consensys/latency/fc-latency-map/map/data/export_${date}.json`;

  try {
    if (fs.existsSync(filepath)) {
      console.error("Exists");
      const data = JSON.parse(fs.readFileSync(filepath, "utf8"));
      res.status(200).json(data);
    }
  } catch (err) {
    console.error("Error", err);
  }

  res.status(404).json({});
};
