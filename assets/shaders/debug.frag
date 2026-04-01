#version 330 core
in vec2 pass_uv;

out vec4 out_color;

uniform sampler2D u_tex;

bool digit(int d, vec2 p) {
	if(p.x > 1. || p.y > 1. || p.x < 0. || p.y < 0.) return false;
	int s = 0;
	float w = 0.1;
	if(p.y >= 1.-w) s = 0;
	else if(p.x >= 1.-w && p.y >= 0.5) s = 1;
	else if(p.x >= 1.-w && p.y < 0.5) s = 2;
	else if(p.y < w) s = 3;
	else if(p.x < w && p.y < 0.5) s = 4;
	else if(p.x < w && p.y >= 0.5) s = 5;
	else if(p.y >= 0.5-w/2. && p.y < 0.5+w/2.) s = 6;
	else return false;

	switch(d) {
		case 0: return s != 6;
		case 1: return s == 1 || s == 2;
		case 2: return s != 5 && s != 2;
		case 3: return s != 4 && s != 5;
		case 4: return s != 0 && s != 3 && s != 4;
		case 5: return s != 1 && s != 4;
		case 6: return s != 1;
		case 7: return s < 3;
		case 8: return true;
		case 9: return s != 4;
	}
}

bool num(int n, vec2 p) {
	int digits = n == 0 ? 1 : int(log(float(n)) / log(10.)) + 1;
    p.x -= 1.;
    p.x *= float(digits)*1.1;
    p.x += 1.05;
	bool r = false;
	do {
		int d = n%10;
		n /= 10;
		r = r || digit(d, p);
		p.x += 1.1;
	} while(n != 0);
	return r;
}

void main()
{   
    vec2 p = vec2(pass_uv.x, 1.-pass_uv.y) * vec2(1.05, 2.15) - vec2(0.025, 0.05);
	ivec2 dim = textureSize(u_tex,0);
	out_color.rgb = vec3(1.) * float(num(dim.x, p - vec2(0., 1.05)) || num(dim.y, p));

	p = pass_uv * 2. - vec2(1.);
	out_color.r += float(max(abs(p.x), abs(p.y)) >= 0.95);

	vec4 c = texture(u_tex, pass_uv);
	out_color.rgb = out_color.rgb * 0.75 + c.rgb * 0.25;
    out_color.a = 1.;
}