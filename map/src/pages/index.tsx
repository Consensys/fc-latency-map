import React, { useState, useEffect } from "react";
import Head from "next/head";
import dynamic from "next/dynamic";
import getConfig from "next/config";
import { promises as fs } from 'fs';

import styles from "@src/styles/Global.module.css";
import Header from "@src/components/Header";

const { publicRuntimeConfig, serverRuntimeConfig } = getConfig();

interface Props {
  data: any;
}

const Home = (props: Props) => {
  const { data } = props;
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
          minHeight: "600px",
        }}
      >
        <MapWithNoSSR data={data} height={size.x} width={size.y} />
      </div>
    </div>
  );
};

export async function getServerSideProps() {
  let data = undefined
  const exports = await fs.readdir(serverRuntimeConfig.path.exportsMeasures);
  if (exports && exports.length > 0) {
    const sorted = exports.sort()
    const latest = sorted[sorted.length - 1]
    console.log("Serving export:", latest)
    data = await fs.readFile(`${serverRuntimeConfig.path.exportsMeasures}/${latest}`, 'utf-8');
    data = JSON.stringify(data);
  }
  return {
    props: {
      data,
    },
  };
}

export default Home;
