(function () {
    'use strict';

    let node_data = {
        elems: [],
        ids: [],
        id2index: {},
        target: null,
        x: 0,
        y: 0,
    };
    const node_event = {
        mouse_down(e) {
            const target = e.target;
            //const targetW = target.offsetWidth;
            const target_rect = target.getBoundingClientRect();

            console.log(target)
            node_data.target = target;
            node_data.x = target_rect.left - e.pageX;
            node_data.y = target_rect.top - e.pageY;
            //target.style.width = `${targetW}px`;
            target.classList.add("moving");
            window.addEventListener("mousemove", node_event.mouse_move);
            window.addEventListener("mouseup", node_event.mouse_up);
        },

        mouse_move(e) {
            const target = node_data.target;
            const x = e.pageX + node_data.x;
            const y = e.pageY + node_data.y;
            target.style.left = `${x}px`;
            target.style.top = `${y}px`;
            console.log(target.style)
            console.log(target.style.left,target.style.top)
            console.log(x,y)
        },

        mouse_up() {
            const target = node_data.target;
            target.classList.remove("moving");
            window.removeEventListener("mousemove", node_event.mouse_move);
            window.removeEventListener("mouseup", node_event.mouse_up);
        }
    };

    window.addEventListener("load",async function(e){
        const res = await fetch("http://localhost:54321/class", {
            method: "GET",
            mode: "cors",
            //cache: "no-cache",
            headers: {
                //"Accept": "application/json",
                "Origin": "http://localhost:8080"
            },
            redirect: "follow",
        });
        console.log(res);
        return;
        for(let i=0;i<nodes.length;i++){
            let node=nodes.item(i);
            node.addEventListener("mousedown",node_event.mouse_down)
        }
    })
}());