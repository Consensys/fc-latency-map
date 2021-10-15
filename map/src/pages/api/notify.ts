import type { NextApiRequest, NextApiResponse } from "next";

export default async (req: NextApiRequest, res: NextApiResponse) => {
  if (req.method === "POST") {
    res.status(200).json({ success: true });
  } else {
    res.status(404).json({});
  }
};
