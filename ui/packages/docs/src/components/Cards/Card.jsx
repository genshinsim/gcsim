import React from "react";
import styled from "styled-components";

const Trunc = styled.p`
    min-width: 0;
    // white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    display: -webkit-box;
    -webkit-line-clamp: 2; /* number of lines to show */
            line-clamp: 2; 
    -webkit-box-orient: vertical;
`


export default function Card({ title, text, link }) {
    return (
        <div class="card-demo">
            <div class="card">
                <div class="card__header">
                    <h3>{title}</h3>
                </div>
                <div class="card__body">
                    <Trunc>
                        {text}
                    </Trunc>
                </div>
                <div class="card__footer">
                    <a class="button button--secondary button--block" href={link}>View</a>
                </div>
            </div>
        </div>
    )
}

