use wasm_bindgen::prelude::*;
use wasm_bindgen::JsCast;
use web_sys::{WebGlProgram, WebGlRenderingContext, WebGlShader, WebGlUniformLocation, WebGlBuffer, console};

#[wasm_bindgen(start)]
pub fn start() -> Result<(), JsValue> {
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
        uniform vec4 u_Color;

        attribute vec3 a_Vertex;

        void main() {
            gl_Position = u_Transform * vec4(a_Vertex, 1.0);
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

        uniform vec4 u_Color;

        void main() {
            gl_FragColor = u_Color;
        }
    "#,
    )?;
    let program = link_program(&context, &vert_shader, &frag_shader)?;
    context.use_program(Some(&program));

    let vertices: [f32; 36] = [0.5, -0.25, 0.25, 0.0, 0.25, 0.0, -0.5, -0.25, 0.25, -0.5, -0.25, 0.25, 0.0, 0.25, 0.0, 0.0, -0.25, -0.5, 0.0, -0.25, -0.5, 0.0, 0.25, 0.0, 0.5, -0.25, 0.25, 0.0, -0.25, -0.5, 0.5, -0.25, 0.25, -0.5, -0.25, 0.25];

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
        let vert_array = js_sys::Float32Array::view(&vertices);

        context.buffer_data_with_array_buffer_view(
            WebGlRenderingContext::ARRAY_BUFFER,
            &vert_array,
            WebGlRenderingContext::STATIC_DRAW,
        );
    }

    let u_color_location     = context.get_uniform_location(&program, "u_Color");
    let u_transform_location = context.get_uniform_location(&program, "u_Transform");
    let a_vertex_location    = context.get_attrib_location(&program,  "a_Vertex");

    render(&context, &[
        0.9968330,  0.0794111,  0.0042227, 0.0,
        -0.0794111,  0.9912028,  0.1058814, 0.0,
        0.0042227, -0.1058814,  0.9943698, 0.0,
        0.0, 0.0, 0.0, 1.0,

        // 2.0, 0.0, 0.0, 0.0,
        // 0.0, 2.0, 0.0, 0.0,
        // 0.0, 0.0, 2.0, 0.0,
        // 0.0, 0.0, 0.0, 1.0,
    ], 
    u_transform_location.as_ref(),
    u_color_location.as_ref(),
    a_vertex_location, 
    Some(&buffer),
    4,
    );

    Ok(())
}

pub fn render(
    context: &WebGlRenderingContext,
    transform: &[f32],
    u_transform_location: Option<&WebGlUniformLocation>,
    u_color_location: Option<&WebGlUniformLocation>,
    a_vertex_location: i32,
    buffer: Option<&WebGlBuffer>,
    triangles: i32,
) {
    context.clear_color(0.0, 0.5, 0.0, 1.0);
    context.clear(WebGlRenderingContext::COLOR_BUFFER_BIT);

    context.uniform_matrix4fv_with_f32_array(u_transform_location, false, transform);
    context.uniform4fv_with_f32_array(u_color_location, &[1.0, 0.0, 0.0, 1.0]);
    context.bind_buffer(WebGlRenderingContext::ARRAY_BUFFER, buffer);
    context.vertex_attrib_pointer_with_i32(a_vertex_location as u32, 3, WebGlRenderingContext::FLOAT, false, 0, 0);
    context.enable_vertex_attrib_array(a_vertex_location as u32);

    context.draw_arrays(
        WebGlRenderingContext::TRIANGLES,
        0,
        triangles * 3
    );
    console::log_1(&a_vertex_location.into());

    context.uniform4fv_with_f32_array(u_color_location, &[0.0, 0.0, 0.0, 1.0]);
    for start in 0..triangles {
        context.draw_arrays(
            WebGlRenderingContext::LINE_LOOP,
            start*3,
            3
        );
    }
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