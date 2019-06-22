
const genTree = (nodes) => {
    console.log(nodes)
}

config = {
    container: "#tree"
};
parent_node = {
    text: { name: "Parent node" }
};
first_child = {
    parent: parent_node,
    text: { name: "First child" }
};
second_child = {
    parent: parent_node,
    text: { name: "Second child" }
};
second_first_child = {
    parent: first_child,
    text: { name: "Second First child" }
};
simple_chart_config = [
    config, parent_node,
    first_child, second_child, second_first_child
];
let chart = new Treant(simple_chart_config, function() { alert( 'Tree Loaded' ) }, $ );