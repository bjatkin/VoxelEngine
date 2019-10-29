use web_sys::{WebGlBuffer, WebGlRenderingContext, WebGlProgram, WebGlShader};

pub struct WebGL {}

impl WebGL {
    pub fn compile_shader(
        gl: &WebGlRenderingContext,
        shader_type: u32,
        source: &str,
    ) -> Result<WebGlShader, String> {
        let shader = gl
            .create_shader(shader_type)
            .ok_or_else(|| String::from("Unable to create shader object"))?;
        gl.shader_source(&shader, source);
        gl.compile_shader(&shader);

        if gl
            .get_shader_parameter(&shader, WebGlRenderingContext::COMPILE_STATUS)
            .as_bool()
            .unwrap_or(false)
        {
            Ok(shader)
        } else {
            Err(gl
                .get_shader_info_log(&shader)
                .unwrap_or_else(|| String::from("Unknown error creating shader")))
        }
    }

    pub fn link_program(
        gl: &WebGlRenderingContext,
        vert_shader: &WebGlShader,
        frag_shader: &WebGlShader,
    ) -> Result<WebGlProgram, String> {
        let program = gl
            .create_program()
            .ok_or_else(|| String::from("Unable to create shader object"))?;

        gl.attach_shader(&program, vert_shader);
        gl.attach_shader(&program, frag_shader);
        gl.link_program(&program);

        if gl
            .get_program_parameter(&program, WebGlRenderingContext::LINK_STATUS)
            .as_bool()
            .unwrap_or(false)
        {
            Ok(program)
        } else {
            Err(gl
                .get_program_info_log(&program)
                .unwrap_or_else(|| String::from("Unknown error creating program object")))
        }
    }

    // pub fn create_buffer(
    //     gl: &WebGlRenderingContext,
    //     data: &[f32],
    // ) -> Result<WebGlBuffer, String> {
    //     let buffer = gl.create_buffer().ok_or("failed to create buffer")?;
    //     gl.bind_buffer(WebGlRenderingContext::ARRAY_BUFFER, Some(&buffer));

    //     // Note that `Float32Array::view` is somewhat dangerous (hence the
    //     // `unsafe`!). This is creating a raw view into our module's
    //     // `WebAssembly.Memory` buffer, but if we allocate more pages for ourself
    //     // (aka do a memory allocation in Rust) it'll cause the buffer to change,
    //     // causing the `Float32Array` to be invalid.
    //     //
    //     // As a result, after `Float32Array::view` we have to be very careful not to
    //     // do any memory allocations before it's dropped.
    //     unsafe {
    //         let data_array = js_sys::Float32Array::view(&data);

    //         gl.buffer_data_with_array_buffer_view(
    //             WebGlRenderingContext::ARRAY_BUFFER,
    //             &data_array,
    //             WebGlRenderingContext::STATIC_DRAW,
    //         );
    //     }

    //     Ok(buffer)
    // }

    pub fn update_buffer<'a>(
        gl: &WebGlRenderingContext,
        data: &'a [f32],
        buffer: &'a WebGlBuffer,
    ) -> &'a WebGlBuffer {
        gl.bind_buffer(WebGlRenderingContext::ARRAY_BUFFER, Some(&buffer));

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

            gl.buffer_data_with_array_buffer_view(
                WebGlRenderingContext::ARRAY_BUFFER,
                &data_array,
                WebGlRenderingContext::STATIC_DRAW,
            );
        }

        buffer
    }
}