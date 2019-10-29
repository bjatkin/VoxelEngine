use wasm_bindgen::prelude::*;
use wasm_bindgen::JsCast;
use web_sys::{WebGlProgram, WebGlRenderingContext, WebGlShader, WebGlUniformLocation, WebGlBuffer, console};
use std::cell::RefCell;
use std::rc::Rc;
mod vector;
mod voxel;
mod scene;
mod web_gl;

#[wasm_bindgen(start)]
pub fn run() -> Result<(), JsValue> {
    let document = web_sys::window().unwrap().document().unwrap();
    let canvas = document.get_element_by_id("canvas").unwrap();
    let canvas: web_sys::HtmlCanvasElement = canvas.dyn_into::<web_sys::HtmlCanvasElement>()?;

    let gl = canvas
        .get_context("webgl")?
        .unwrap()
        .dyn_into::<WebGlRenderingContext>()?;
    
    let scene = scene::Scene::new(&gl)?;

    let test_voxel = voxel::Voxel::new(
        vector::Vector::new(-0.5, -0.25, 0.25),
        [
            vector::Vector::new(1.0, 0.0, 0.0),
            vector::Vector::new(0.0, 1.0, 0.0),
            vector::Vector::new(0.0, 0.0, 1.0),
            vector::Vector::new(1.0, 0.0, 1.0),
            vector::Vector::new(1.0, 1.0, 0.0),
            vector::Vector::new(1.0, 1.0, 1.0),
        ]
    );

    scene.add_voxel(test_voxel);

    //Prep the closure for request animiation frame
    let f = Rc::new(RefCell::new(None));
    let g = f.clone();

    *g.borrow_mut() = Some(Closure::wrap(Box::new(move || {
        scene.render();
        request_animation_frame(f.borrow().as_ref().unwrap());
    }) as Box<dyn FnMut()>));

    request_animation_frame(g.borrow().as_ref().unwrap());
    Ok(())
}

pub fn oldRun() -> Result<(), JsValue> {
    let document = web_sys::window().unwrap().document().unwrap();
    let canvas = document.get_element_by_id("canvas").unwrap();
    let canvas: web_sys::HtmlCanvasElement = canvas.dyn_into::<web_sys::HtmlCanvasElement>()?;

    let context = canvas
        .get_context("webgl")?
        .unwrap()
        .dyn_into::<WebGlRenderingContext>()?;

    let vert_shader = compile_shader(
        &context,
        WebGlRenderingContext::VERTEX_SHADER,
        r#"
        // VertexShader

        precision mediump int;
        precision mediump float;

        uniform mat4 u_Transform;

        attribute vec3 a_Vertex;
        attribute vec3 a_Color;

        varying vec4 v_vertex_color;

        void main() {
            gl_Position = u_Transform * vec4(a_Vertex, 1.0);

            v_vertex_color = vec4(a_Color, 1.0);
        }
        "#,
    )?;
    let frag_shader = compile_shader(
        &context,
        WebGlRenderingContext::FRAGMENT_SHADER,
        r#"
        // Fragment shader

        precision mediump int;
        precision mediump float;

        varying vec4 v_vertex_color;

        void main() {
            gl_FragColor = v_vertex_color;
        }
    "#,
    )?;
    let program = link_program(&context, &vert_shader, &frag_shader)?;
    context.use_program(Some(&program));

    let _ = voxel::Voxel::new(
        vector::Vector::new(-0.5, -0.25, 0.25),
        [
            vector::Vector::new(1.0, 0.0, 0.0),
            vector::Vector::new(1.0, 0.0, 0.0),
            vector::Vector::new(1.0, 0.0, 0.0),
            vector::Vector::new(1.0, 0.0, 0.0),
            vector::Vector::new(1.0, 0.0, 0.0),
            vector::Vector::new(1.0, 0.0, 0.0),
        ]
    );

    let (vertices, colors) = create_cube(-0.5, -0.25, 0.25, 0.8, 0.0, 0.0);
    // let vertices: [f32; 36] = [
    //     //triangle 1
    //     0.5, -0.25, 0.25, 
    //     0.0, 0.25, 0.0, 
    //     -0.5, -0.25, 0.25,
    //     //triangle 2
    //     -0.5, -0.25, 0.25,
    //     0.0, 0.25, 0.0,
    //     0.0, -0.25, -0.5, 
    //     //triangle 3 
    //     0.0, -0.25, -0.5, 
    //     0.0, 0.25, 0.0, 
    //     0.5, -0.25, 0.25,
    //     //triangle 4 
    //     0.0, -0.25, -0.5,
    //     0.5, -0.25, 0.25,
    //     -0.5, -0.25, 0.25
    //     ];

    // let colors: [f32; 36] = [
    //     //triangle 1
    //     0.0, 0.0, 1.0,
    //     0.0, 1.0, 0.0,
    //     1.0, 0.0, 1.0,
    //     //triangle 2
    //     1.0, 0.0, 1.0,
    //     0.0, 1.0, 0.0,
    //     1.0, 0.0, 0.0,
    //     //triangle 3
    //     1.0, 0.0, 0.0,
    //     0.0, 1.0, 0.0,
    //     0.0, 0.0, 1.0,
    //     //triangle 4
    //     0.0, 0.0, 1.0,
    //     1.0, 0.0, 1.0,
    //     1.0, 0.0, 0.0,
    //     ];

    // let buffer = context.create_buffer().ok_or("failed to create buffer")?;
    // context.bind_buffer(WebGlRenderingContext::ARRAY_BUFFER, Some(&buffer));

    // // Note that `Float32Array::view` is somewhat dangerous (hence the
    // // `unsafe`!). This is creating a raw view into our module's
    // // `WebAssembly.Memory` buffer, but if we allocate more pages for ourself
    // // (aka do a memory allocation in Rust) it'll cause the buffer to change,
    // // causing the `Float32Array` to be invalid.
    // //
    // // As a result, after `Float32Array::view` we have to be very careful not to
    // // do any memory allocations before it's dropped.
    // unsafe {
    //     let vert_array = js_sys::Float32Array::view(&vertices);

    //     context.buffer_data_with_array_buffer_view(
    //         WebGlRenderingContext::ARRAY_BUFFER,
    //         &vert_array,
    //         WebGlRenderingContext::STATIC_DRAW,
    //     );
    // }
    let triangles_vertex_buffer_id = create_buffer(&context, &vertices)?;
    let triangles_color_buffer_id = create_buffer(&context, &colors)?;

    let u_transform_location = context.get_uniform_location(&program, "u_Transform");

    let u_color_location     = context.get_attrib_location(&program, "a_Color");
    let a_vertex_location    = context.get_attrib_location(&program,  "a_Vertex");

    //Prep the closure for request animiation frame
    let f = Rc::new(RefCell::new(None));
    let g = f.clone();

    let mut theta: f32 = 1.0;
    let step: f32 = std::f32::consts::PI/48.0;
    *g.borrow_mut() = Some(Closure::wrap(Box::new(move || {
        theta += step;
        if theta >= 2.0*std::f32::consts::PI {
            theta = 0.0;
        }

        render(&context, &[
             theta.cos(), -theta.sin(), 0.0, 0.0,
            theta.sin(), theta.cos(), 0.0, 0.0,
            0.0, 0.0, 1.0, 0.0,
            0.0, 0.0, 0.0, 1.0,
        ], 
        u_transform_location.as_ref(),
        u_color_location,
        a_vertex_location, 
        Some(&triangles_vertex_buffer_id),
        Some(&triangles_color_buffer_id),
        4,
        );

        request_animation_frame(f.borrow().as_ref().unwrap());
    }) as Box<dyn FnMut()>));

    request_animation_frame(g.borrow().as_ref().unwrap());
    Ok(())
}

pub fn create_buffer(
    context: &WebGlRenderingContext,
    data: &[f32],
) -> Result<WebGlBuffer, String> {
    let buffer = context.create_buffer().ok_or("failed to create buffer")?;
    context.bind_buffer(WebGlRenderingContext::ARRAY_BUFFER, Some(&buffer));

    // Note that `Float32Array::view` is somewhat dangerous (hence the
    // `unsafe`!). This is creating a raw view into our module's
    // `WebAssembly.Memory` buffer, but if we allocate more pages for ourself
    // (aka do a memory allocation in Rust) it'll cause the buffer to change,
    // causing the `Float32Array` to be invalid.
    //
    // As a result, after `Float32Array::view` we have to be very careful not to
    // do any memory allocations before it's dropped.
    unsafe {
        let data_array = js_sys::Float32Array::view(&data);

        context.buffer_data_with_array_buffer_view(
            WebGlRenderingContext::ARRAY_BUFFER,
            &data_array,
            WebGlRenderingContext::STATIC_DRAW,
        );
    }

    Ok(buffer)
}

pub fn render(
    context: &WebGlRenderingContext,
    transform: &[f32],
    u_transform_location: Option<&WebGlUniformLocation>,
    a_color_location: i32,
    a_vertex_location: i32,
    tri_buffer: Option<&WebGlBuffer>,
    color_buffer: Option<&WebGlBuffer>,
    triangles: i32,
) {
    context.uniform_matrix4fv_with_f32_array(u_transform_location, false, transform);

    context.bind_buffer(WebGlRenderingContext::ARRAY_BUFFER, tri_buffer);

    context.vertex_attrib_pointer_with_i32(a_vertex_location as u32, 3, WebGlRenderingContext::FLOAT, false, 0, 0);
    context.enable_vertex_attrib_array(a_vertex_location as u32);

    context.bind_buffer(WebGlRenderingContext::ARRAY_BUFFER, color_buffer);

    context.vertex_attrib_pointer_with_i32(a_color_location as u32, 3, WebGlRenderingContext::FLOAT, false, 0, 0);
    context.enable_vertex_attrib_array(a_color_location as u32);

    context.draw_arrays(
        WebGlRenderingContext::TRIANGLES,
        0,
        triangles * 3
    );
}

fn window() -> web_sys::Window {
    web_sys::window().expect("no global `window` exists")
}

fn request_animation_frame(f: &Closure<dyn FnMut()>) {
    window()
        .request_animation_frame(f.as_ref().unchecked_ref())
        .expect("should register `requestAnimationFrame` OK");
}

pub fn compile_shader(
    context: &WebGlRenderingContext,
    shader_type: u32,
    source: &str,
) -> Result<WebGlShader, String> {
    let shader = context
        .create_shader(shader_type)
        .ok_or_else(|| String::from("Unable to create shader object"))?;
    context.shader_source(&shader, source);
    context.compile_shader(&shader);

    if context
        .get_shader_parameter(&shader, WebGlRenderingContext::COMPILE_STATUS)
        .as_bool()
        .unwrap_or(false)
    {
        Ok(shader)
    } else {
        console::log_1(&"There was an error creating the shader".into());
        Err(context
            .get_shader_info_log(&shader)
            .unwrap_or_else(|| String::from("Unknown error creating shader")))
    }
}

pub fn link_program(
    context: &WebGlRenderingContext,
    vert_shader: &WebGlShader,
    frag_shader: &WebGlShader,
) -> Result<WebGlProgram, String> {
    let program = context
        .create_program()
        .ok_or_else(|| String::from("Unable to create shader object"))?;

    context.attach_shader(&program, vert_shader);
    context.attach_shader(&program, frag_shader);
    context.link_program(&program);

    if context
        .get_program_parameter(&program, WebGlRenderingContext::LINK_STATUS)
        .as_bool()
        .unwrap_or(false)
    {
        Ok(program)
    } else {
        console::log_1(&"There was an error creating the program".into());
        Err(context
            .get_program_info_log(&program)
            .unwrap_or_else(|| String::from("Unknown error creating program object")))
    }
}

fn create_cube(x: f32, y: f32, z: f32, r: f32, g: f32, b: f32) -> ([f32; 108], [f32; 108]) {
    let verts: [f32; 108] = [
        x+1.0,     y,     z,
        x+1.0, y+1.0,     z,
            x,     y,     z,//f1
            x,     y,     z,
            x, y+1.0,     z,
        x+1.0, y+1.0,     z,//f2
            x, y+1.0,     z,
            x, y+1.0, z-1.0,
        x+1.0, y+1.0, z-1.0,//f3
            x, y+1.0,     z,
        x+1.0, y+1.0,     z,
        x+1.0, y+1.0, z-1.0,//f4
        x+1.0,     y,     z,
        x+1.0, y+1.0,     z,
        x+1.0, y+1.0, z-1.0,//f5
        x+1.0,     y,     z,
        x+1.0, y+1.0, z-1.0,
        x+1.0,     y, z-1.0,//f6
            x,     y,     z,
            x, y+1.0,     z,
            x, y+1.0, z-1.0,//f7
            x,     y,     z,
            x, y+1.0, z-1.0,
            x,     y, z-1.0,//f8
            x,     y,     z,
            x,     y, z-1.0,
        x+1.0,     y, z-1.0,//f9
            x,     y,     z,
        x+1.0,     y, z-1.0,
        x+1.0,     y,     z,//f10
            x,     y, z-1.0,
            x, y+1.0, z-1.0,
        x+1.0, y+1.0, z-1.0,//f11
            x,     y, z-1.0,
        x+1.0, y+1.0, z-1.0,
        x+1.0,     y, z-1.0,//f12
    ];

    let colors: [f32; 108] = [
        r, g, b,
        r, g, b,
        r, g, b,//f1
        r, g, b,
        r, g, b,
        r, g, b,//f2
        r, g, b,
        r, g, b,
        r, g, b,//f3
        r, g, b,
        r, g, b,
        r, g, b,//f4
        r, g, b,
        r, g, b,
        r, g, b,//f5
        r, g, b,
        r, g, b,
        r, g, b,//f6
        r, g, b,
        r, g, b,
        r, g, b,//f7
        r, g, b,
        r, g, b,
        r, g, b,//f8
        r, g, b,
        r, g, b,
        r, g, b,//f9
        r, g, b,
        r, g, b,
        r, g, b,//f10
        r, g, b,
        r, g, b,
        r, g, b,//f11
        r, g, b,
        r, g, b,
        r, g, b,//f12
    ];

    return (verts, colors);
}