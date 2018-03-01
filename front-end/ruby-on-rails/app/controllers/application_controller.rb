class ApplicationController < ActionController::Base
  protect_from_forgery with: :exception

  def sort_column(array, direction, index)
    if direction == "asc"
      array = array.sort_by! { |e| e[index] }
    else 
      array = array.sort_by! { |e| e[index] }.reverse
    end
  end

  def strip_underscore(text)
		return (text.include? "_") ? text.gsub(/_/, ' ') : text
	end

end
