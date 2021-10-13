import type { NextApiResponse } from "next";

export default async ({}, res: NextApiResponse) => {
  res.status(200).json({ success: true });
};
