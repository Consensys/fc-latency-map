import React, { useState, useEffect } from "react";
import Head from "next/head";
import dynamic from "next/dynamic";
import getConfig from "next/config";
import moment, { Moment } from "moment";
import * as fs from "fs";
import * as path from "path";

import styles from "@src/styles/Global.module.css";
import Header from "@src/components/Header";

const { publicRuntimeConfig, serverRuntimeConfig } = getConfig();

interface Props {
  data: any;
  date: string;
  dates: string[];
}

const Home = (props: Props) => {
  const { data, date, dates } = props;
  const MapWithNoSSR = dynamic(() => import("@src/components/map/Map"), {
    ssr: false,
  });

  const [size, setSize] = useState({
    x: 0,
    y: 0,
  });

  const updateSize = () => {
    if (size.x != window.innerWidth || size.y != window.innerHeight) {
      setSize({
        x: window.innerWidth,
        y: window.innerHeight,
      });
    }
  };

  useEffect(() => {
    setSize({
      x: window.innerWidth,
      y: window.innerHeight,
    });
  }, []);

  useEffect(() => (window.onresize = updateSize), []);

  return (
    <div className={styles.container}>
      <Head>
        <title>Filecoin Latency Map</title>
        <meta name="description" content={publicRuntimeConfig.app.name} />
        <link rel="icon" href="/favicon.ico" />
      </Head>
      <Header />
      <div
        id="map"
        style={{
          width: "100%",
          minHeight: "500px",
        }}
      >
        <MapWithNoSSR
          data={data}
          date={date}
          dates={dates}
          height={size.x}
          width={size.y}
        />
      </div>
    </div>
  );
};

export async function getServerSideProps() {
  let data = JSON.stringify({
    locations: [],
    miners: [],
  });
  let date = "";
  let dates: string[] = [];

  const files = fs.readdirSync(serverRuntimeConfig.path.exportsMeasures);
  if (files && files.length > 0) {
    const filesSorted = files.sort();
    const latest = filesSorted[filesSorted.length - 1];

    const datesFound = latest.match(/\d+-\d+-\d+/g);
    if (datesFound) {
      date = datesFound[0];
    }
    const filename = path.join(
      serverRuntimeConfig.path.exportsMeasures,
      latest
    );

    data = await fs.readFileSync(filename, "utf-8");
    dates = files.map((file) => {
      const datesFound = file.match(/\d+-\d+-\d+/g);
      if (datesFound) {
        return datesFound[0];
      }
      return "error";
    });
  }

  return {
    props: {
      dates,
      data,
      date,
    },
  };
}

export default Home;
