#version 330 core

in vec2 pass_uv;

out vec4 out_color;

uniform sampler2DArray s_atlas;
uniform int u_type;

const float smoothing = 1.0/16.0;

float median(float r, float g, float b) {
    return max(min(r, g), min(max(r, g), b));
}

void main() {
	vec3 uvw = vec3(fract(pass_uv.x), pass_uv.y, floor(pass_uv.x));

	float sd;
	float alpha;

	switch(u_type) {
		case(0):
			out_color = texture(s_atlas, uvw);
		break;
		case(1): 
			sd = texture(s_atlas, uvw).a-0.5;
    		alpha = smoothstep(- smoothing, smoothing, sd);
    		out_color = vec4(vec3(1.),  alpha);
		break;
		case(2):
			vec3 msd = texture(s_atlas, uvw).rgb;
			sd = median(msd.x, msd.y, msd.z) - 0.5;
			alpha = smoothstep(- smoothing, smoothing, sd);
			float d = fwidth(sd);
    		alpha = smoothstep(-d, d, sd);
			out_color = vec4(vec3(1.),  alpha);
		break;
	}
}