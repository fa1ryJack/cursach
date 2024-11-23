<script setup>
import { ref, watchEffect } from "vue";
import * as d3 from "d3";

const { likesData } = defineProps(["likesData"]);
const chart = ref(null);

watchEffect(() => {
  setTimeout(() => {
    drawChart();
  }, 10);
});

function drawChart() {
  const grey = "#e5e5e5";
  const bley = "#333333";
  const trDelay = 500;

  // Specify the chart’s dimensions.
  const width = 1000;
  const height = width;

  // Compute the layout.
  const pack = (data) =>
    d3.pack().size([width, height]).padding(3)(
      d3
        .hierarchy(data)
        .sum((d) => d.value)
        .sort((a, b) => b.value - a.value)
    );
  const root = pack(likesData);

  d3.select("svg").selectAll("*").remove();

  // Create the SVG container.
  const svg = d3
    .select("svg")
    .attr("viewBox", `-${width / 2} -${height / 2} ${width} ${height}`)
    .attr("width", width)
    .attr("height", height)
    .attr(
      "style",
      `max-width: 100%; height: auto; display: block; margin: 0 -14px; background: #090909; cursor: pointer;`
    );

  // Create defs for storing images
  const defs = svg.append("defs");
  root.children.forEach(function (d) {
    defs
      .append("pattern")
      .attr("id", "id" + d.data.id)
      .attr("patternUnits", "userSpaceOnUse")
      .append("svg:image")
      .attr("xlink:href", d.data.AvatarURL);
    d.children.forEach(function (d) {
      defs
        .append("pattern")
        .attr("id", "id" + d.data.id)
        .attr("patternUnits", "userSpaceOnUse")
        .append("svg:image")
        .attr("xlink:href", d.data.ArtworkURL);
    });
  });

  // Append the nodes
  const node = svg
    .append("g")
    .selectAll("circle")
    .data(root.descendants().slice(1))
    .join("circle")
    .style("fill", (d) => `url(#${"id" + d.data.id})`)
    .attr("stroke-width", "3px")
    .on("mouseover", function () {
      d3.select(this).attr("stroke", "#ff5500");
    })
    .on("mouseout", function () {
      d3.select(this).attr("stroke", null);
    })
    .on("click", function (event, d) {
      focus !== d && (zoom(event, d), event.stopPropagation());
    });

  // Append the text labels.
  const label = svg
    .append("g")
    .style("font", "20px sans-serif")
    .style("fill", "#ff5500")
    .attr("pointer-events", "none")
    .attr("text-anchor", "middle")
    .selectAll("text")
    .data(root.descendants())
    .join("text")
    .style("fill-opacity", (d) => (d.parent === root ? 1 : 0))
    .style("display", (d) => (d.parent === root ? "inline" : "none"))
    .text((d) => d.data.name);

  // Create the zoom behavior and zoom immediately in to the initial focus node.
  svg.on("click", (event) => zoom(event, root));
  let focus = root;
  let view;
  zoomTo([focus.x, focus.y, focus.r * 2]);

  function zoomTo(v) {
    const k = width / v[2];
    view = v;
    label.attr(
      "transform",
      (d) => `translate(${(d.x - v[0]) * k},${(d.y - v[1]) * k})`
    );
    node.attr(
      "transform",
      (d) => `translate(${(d.x - v[0]) * k},${(d.y - v[1]) * k})`
    );
    node.attr("r", (d) => d.r * k);

    root.children.forEach(function (d) {
      defs
        .select(`#${"id" + d.data.id}`)
        .attr("width", d.r * 2 * k)
        .attr("height", d.r * 2 * k)
        .attr("x", d.r * k)
        .attr("y", d.r * k)
        .select("image")
        .attr("width", d.r * 2 * k)
        .attr("height", d.r * 2 * k);

      d.children.forEach(function (d) {
        defs
          .select(`#${"id" + d.data.id}`)
          .attr("width", d.r * 2 * k)
          .attr("height", d.r * 2 * k)
          .attr("x", d.r * k)
          .attr("y", d.r * k)
          .select("image")
          .attr("width", d.r * 2 * k)
          .attr("height", d.r * 2 * k);
      });
    });
  }

  function zoom(event, d) {
    const focus0 = focus;

    focus = d;

    if (focus0.depth === 2 && focus.depth === 0) {
      console.log("Targeted song:", focus0.data.name, focus0.data.id);
    }

    const transition = svg
      .transition()
      // .delay(trDelay * 2)
      .duration(event.altKey ? 7500 : 750)
      .tween("zoom", (d) => {
        const i = d3.interpolateZoom(view, [focus.x, focus.y, focus.r * 2]);
        return (t) => zoomTo(i(t));
      });

    function changeFillCondition(d, focus) {
      return (
        d.data.name === focus.data.name || //current d is target
        (d.data.name === focus?.parent?.data?.name && !focus.children) || //current d is parent of target and target is song
        (d.parent.data.name === focus?.parent?.data?.name &&
          !focus.children &&
          !d.children) || //current d is a track sibling of target track
        (!d.children && d.parent.data.name === focus.data.name) //target is parent of current d track
      );
    }

    function relativesCondition(focus0, focus) {
      if (focus.depth !== 0 && focus0.depth === 0) {
        return true;
      }
      if (focus.depth === 0) {
        return false;
      }
      if (
        (focus.parent === focus0.parent && focus.parent.depth !== 0) ||
        focus.parent === focus0 ||
        focus0.parent === focus
      ) {
        return false;
      }
      return true;
    }

    //-----changing fill-----
    if (relativesCondition(focus0, focus)) {
      node
        .transition()
        .duration(trDelay)
        .style("opacity", function (d) {
          if (changeFillCondition(d, focus)) {
            return 1;
          } else {
            return 0;
          }
        })
        .style("fill", function (d) {
          if (changeFillCondition(d, focus)) {
            return `url(#${"id" + d.data.id})`;
          }
          if (d.children) {
            return bley;
          } else {
            return grey;
          }
        })
        .transition()
        .duration(trDelay)
        .style("opacity", 1);
    }

    if (focus.depth === 0) {
      setTimeout(function () {
        node
          .transition()
          .duration(trDelay)
          .style("opacity", 0)
          .style("fill", (d) => `url(#${"id" + d.data.id})`)
          .transition()
          .delay(trDelay)
          .duration(trDelay)
          .style("opacity", 1);
      }, trDelay * 2);
    }
    //-----changing fill-----

    if ((focus.depth !== 0 && focus0.depth !== 0) || focus.depth !== 0) {
      label
        .transition(transition)
        .delay(trDelay)
        .style("font", "30px Helvetica");
    } else {
      label
        .transition(transition)
        .delay(trDelay)
        .style("font", "20px Helvetica");
    }
    label
      .transition(transition)
      .delay(trDelay)
      .style("fill-opacity", (d) =>
        d.parent === focus || (!d.children && focus === d) ? 1 : 0
      )
      .on("start", function (d) {
        if (d.parent === focus || (!d.children && focus === d))
          this.style.display = "inline";
      });
  }
  console.log(svg.node());
}
</script>

<template>
  <div>
    <h3>
      <span>{{ likesData.name }}</span
      >'s likes grouped by uploaders:
    </h3>
    <svg></svg>
  </div>
</template>

<!--   <iframe
    width="80%"
    height="300"
    src="https://w.soundcloud.com/player/?url=https%3A//api.soundcloud.com/tracks/81345014&color=%23ff5500&auto_play=false&hide_related=false&show_comments=true&show_user=true&show_reposts=false&show_teaser=true&visual=true"
  ></iframe> -->

<style lang="scss" scoped>
@use "../assets/colors";
span {
  color: colors.$orange;
}
</style>
