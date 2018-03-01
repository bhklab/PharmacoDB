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

	def sortable(column, title = nil, anchor = nil)
		title ||= column.titleize
		direction = column == params[:sort] && params[:direction] == "asc" ? "desc" : "asc"
		if direction === "asc"
			content_tag(:div, link_to( (title + "&nbsp;" + content_tag(:i, nil, :class => "ion-arrow-up-b")).html_safe , {:sort => column, :direction => direction, :anchor => anchor}))
		else 
			content_tag(:div, link_to( (title + "&nbsp;" + content_tag(:i, nil, :class => "ion-arrow-down-b")).html_safe , {:sort => column, :direction => direction, :anchor => anchor}))
		end
	end

	# Strip all whitespace between the HTML tags in the passed block, and
	# on its start and end.
	def spaceless(&block)
		contents = capture(&block)
		contents.strip.gsub(/>\s+</, '><').html_safe
	end

end
