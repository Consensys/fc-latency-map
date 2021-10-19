import { Card, Row, Col, Statistic } from "antd";
import getConfig from "next/config";
import * as am4core from "@amcharts/amcharts4/core";
import am4themes_animated from "@amcharts/amcharts4/themes/animated";

import styles from "@src/styles/Global.module.css";

am4core.useTheme(am4themes_animated);

const { publicRuntimeConfig } = getConfig();

interface Props {
  miner: any;
}

const Miner = (props: Props) => {
  const { miner } = props;

  const displayMiner = (miner: any) => {
    var color = "#000000";
    if (miner.latency && miner.latency.avg) {
      if (miner.latency.avg < 0) {
        color = "#000000";
      } else if (miner.latency.avg < publicRuntimeConfig.latency.low) {
        color = "#00ff00";
      } else if (miner.latency.avg < publicRuntimeConfig.latency.medium) {
        color = "#FFFF00";
      } else {
        color = "#ff5d54";
      }
    }

    return (
      <Card title="Miner" className={styles.information}>
        <Row gutter={[16, 16]}>
          <Col span={12}>
            <p>
              Address: <strong>{miner.title}</strong>
            </p>
            <p>IP: {miner.ip}</p>
            <p>Latitude: {miner.latitude}</p>
            <p>Longitude: {miner.longitude}</p>
            {/* <p>
          Latency:
          {miner.latency ? <>{miner.latency.avg}</> : <> N/A</>}
        </p> */}
          </Col>
          <Col span={12}>
            {/* <p> */}
            {/* Latency:
          {miner.latency ? <>{miner.latency.avg}</> : <> N/A</>} */}
            <Statistic
              title="Latency"
              value={
                miner.latency && miner.latency.avg > 0
                  ? miner.latency.avg
                  : "N/A"
              }
              precision={2}
              // valueStyle={{ color }}
              style={{
                borderColor: color,
              }}
              className={styles.latencyStat}
              // prefix={<ArrowUpOutlined />}
              // suffix="%"
            />
            {/* </p> */}
          </Col>
        </Row>
      </Card>
    );
  };

  return <div>{miner && displayMiner(miner)}</div>;
};

export default Miner;
