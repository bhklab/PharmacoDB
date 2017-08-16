module ApplicationHelper

	# Returns the full title on a per-page basis.
	def full_title(page_title = '')
		base_title = "PharmacoDB"
		if page_title.empty?
			base_title
		else
			page_title + " | " + base_title
		end
	end

	# Strip all whitespace between the HTML tags in the passed block, and
	# on its start and end.
	def spaceless(&block)
		contents = capture(&block)
		contents.strip.gsub(/>\s+</, '><').html_safe
	end

end
