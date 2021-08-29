import React from "react";

import {
  Button,
  ButtonGroup,
  Card,
  Elevation,
  H4,
  Icon,
} from "@blueprintjs/core";

function TeamBuilder() {
  return (
    <div className="row">
      <div className="col-xs">
        <div className="box">
          <Card style={{ marginTop: "10px" }} elevation={Elevation.THREE}>
            <div className="row">
              <div className="col-xs-2">
                <img
                  src="/assets/UI_AvatarIcon_Ambor.png"
                  alt="amber"
                  style={{ width: "100%" }}
                />
              </div>
              <div className="col-xs-10">
                <div className="row">
                  <div className="col-xs">
                    <H4>Basic</H4>
                    <Card elevation={Elevation.TWO}>
                      <div>
                        <div className="row">
                          <span className="col-xs">Level</span>
                          <span
                            className="col-xs"
                            style={{ textAlign: "right" }}
                          >
                            1234
                          </span>
                        </div>
                        <div className="row">
                          <span className="col-xs">Ascension</span>
                          <span
                            className="col-xs"
                            style={{ textAlign: "right" }}
                          >
                            1234
                          </span>
                        </div>
                        <div className="row">
                          <span className="col-xs">Cons</span>
                          <span
                            className="col-xs"
                            style={{ textAlign: "right" }}
                          >
                            1234
                          </span>
                        </div>
                        <div className="row">
                          <span className="col-xs">Auto</span>
                          <span
                            className="col-xs"
                            style={{ textAlign: "right" }}
                          >
                            1234
                          </span>
                        </div>
                        <div className="row">
                          <span className="col-xs">Skill</span>
                          <span
                            className="col-xs"
                            style={{ textAlign: "right" }}
                          >
                            1234
                          </span>
                        </div>
                        <div className="row">
                          <span className="col-xs">Burst</span>
                          <span
                            className="col-xs"
                            style={{ textAlign: "right" }}
                          >
                            1234
                          </span>
                        </div>
                      </div>
                    </Card>
                  </div>
                  <div className="col-xs">
                    <H4>Stats</H4>
                    <Card elevation={Elevation.TWO}>
                      <div>
                        <div className="row">
                          <span className="col-xs">HP</span>
                          <span
                            className="col-xs"
                            style={{ textAlign: "right" }}
                          >
                            1234
                          </span>
                        </div>
                        <div className="row">
                          <span className="col-xs">ATK</span>
                          <span
                            className="col-xs"
                            style={{ textAlign: "right" }}
                          >
                            1234
                          </span>
                        </div>
                        <div className="row">
                          <span className="col-xs">DEF</span>
                          <span
                            className="col-xs"
                            style={{ textAlign: "right" }}
                          >
                            1234
                          </span>
                        </div>
                        <div className="row">
                          <span className="col-xs">EM</span>
                          <span
                            className="col-xs"
                            style={{ textAlign: "right" }}
                          >
                            1234
                          </span>
                        </div>
                        <div className="row">
                          <span className="col-xs">ER</span>
                          <span
                            className="col-xs"
                            style={{ textAlign: "right" }}
                          >
                            1234
                          </span>
                        </div>
                        <div className="row">
                          <span className="col-xs">Crit</span>
                          <span
                            className="col-xs"
                            style={{ textAlign: "right" }}
                          >
                            1234
                          </span>
                        </div>
                        <div className="row">
                          <span className="col-xs">CD</span>
                          <span
                            className="col-xs"
                            style={{ textAlign: "right" }}
                          >
                            1234
                          </span>
                        </div>
                        <div className="row">
                          <span className="col-xs">Electro %</span>
                          <span
                            className="col-xs"
                            style={{ textAlign: "right" }}
                          >
                            1234
                          </span>
                        </div>
                      </div>
                    </Card>
                  </div>
                  <div className="col-xs">
                    <H4>Weapon</H4>
                    <Card elevation={Elevation.TWO}>
                      <div className="row">
                        <span className="col-xs">Amos' Bow</span>
                      </div>
                      <div className="row">
                        <span className="col-xs">Level</span>
                        <span className="col-xs" style={{ textAlign: "right" }}>
                          1234
                        </span>
                      </div>
                      <div className="row">
                        <span className="col-xs">Refine</span>
                        <span className="col-xs" style={{ textAlign: "right" }}>
                          1234
                        </span>
                      </div>
                    </Card>
                  </div>
                  <div className="col-xs">
                    <H4>Artifact Sets</H4>
                    <Card elevation={Elevation.TWO}>
                      <div className="row">
                        <span className="col-xs">Flower</span>
                        <span className="col-xs" style={{ textAlign: "right" }}>
                          1234
                        </span>
                      </div>
                      <div className="row">
                        <span className="col-xs">Feather</span>
                        <span className="col-xs" style={{ textAlign: "right" }}>
                          <Icon icon="edit" />
                        </span>
                      </div>
                      <div className="row">
                        <span className="col-xs">Sands</span>
                        <span className="col-xs" style={{ textAlign: "right" }}>
                          1234
                        </span>
                      </div>
                      <div className="row">
                        <span className="col-xs">Goblet</span>
                        <span className="col-xs" style={{ textAlign: "right" }}>
                          1234
                        </span>
                      </div>
                      <div className="row">
                        <span className="col-xs">Circle</span>
                        <span className="col-xs" style={{ textAlign: "right" }}>
                          1234
                        </span>
                      </div>
                    </Card>
                  </div>
                </div>
              </div>
            </div>
            <ButtonGroup fill style={{ marginTop: "15px" }}>
              <Button>Edit</Button>
              <Button>Clear</Button>
            </ButtonGroup>
          </Card>
        </div>
      </div>
    </div>
  );
}

export default TeamBuilder;
